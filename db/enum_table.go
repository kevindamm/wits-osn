// Copyright (c) 2024 Kevin Damm
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// github:kevindamm/wits-osn/db/enum_table.go

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	osn "github.com/kevindamm/wits-osn"
)

// Serialized enum types are a recurring pattern in fully defined relational DBs
// and here we make the simplifying assumption that any enumerated type can fit
// inside a single byte.  If the domain is larger than 255 items, express that
// relation as a Table[T] and TableIndex[T, comparable] where T is a [Record].
type EnumTable[EnumType ~uint8] struct {
	*Table[EnumRecord[EnumType]]
	TableIndex[EnumRecord[EnumType], string]
	Zero EnumRecord[EnumType]
}

// EnumTable is a common pattern of (id, name) pairwise relations.
// The `id` is its enumeration and the `name` is persisted in the DB for
// documentation.  The list of pairs is provided by a Functor: int -> string.
func MakeEnumTable[EnumType ~uint8](tablename string, enumrange EnumType) EnumTable[EnumType] {
	enumtable := Table[EnumRecord[EnumType]]{
		Name:    tablename,
		Primary: "id",
	}
	enumindex := TableIndex[EnumRecord[EnumType], string]{
		Table:  &enumtable,
		Column: "name",
	}
	return EnumTable[EnumType]{
		&enumtable,
		enumindex,
		EnumRecord[EnumType]{EnumType(0), enumrange},
	}
}

type EnumRecord[EnumType ~uint8] struct {
	// It is assumed that the zero value for the underlying type is UNKNOWN-equiv.
	Value EnumType

	// Represents the limit value (excl.) for the valid enumeration range.
	Range EnumType
}

func (EnumRecord[T]) SqlCreate(tablename string) string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "name"  VARCHAR(9) NOT NULL
  ) WITHOUT ROWID;`, tablename)
}

func (record EnumRecord[T]) SqlInit(tablename string) string {
	values := make([]string, osn.FetchStatusRange)
	for i := range record.Range {
		status := T(i)
		values[i] = fmt.Sprintf(`(%d, "%s")`, status, status)
	}
	return fmt.Sprintf(`INSERT INTO %s VALUES %s;`,
		tablename, strings.Join(values, ", "))
}

func (record EnumRecord[T]) IsValid() bool {
	return record.Value < record.Range
}

func (record EnumRecord[T]) ScanRecord(row *sql.Row) error {
	// The underlying uint value can easily be scanned by database/sql.
	return row.Scan(record)
}

func (record EnumRecord[T]) RecordValues() ([]driver.Value, error) {
	return []driver.Value{record.Value, fmt.Sprintf("%s", record.Value)}, nil
}
