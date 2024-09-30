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
	"errors"
	"fmt"
	"strings"

	osn "github.com/kevindamm/wits-osn"
)

// Serialized enum types are a recurring pattern in fully defined relational DBs
// and here we make the simplifying assumption that any enumerated type can fit
// inside a single byte.  If the domain is larger than 255 items, express that
// relation as a [Table[T Record]], named by a [TableIndex[Record, comparable]].
type EnumTable[T osn.EnumType] struct {
	table[EnumRecord[T]]

	values []T
	naming map[string]T
}

// Constructor function for making an enum table with the indicated values.
func MakeEnumTable[T osn.EnumType](tablename string, values []T) EnumTable[T] {
	naming := make(map[string]T)
	for _, t := range values {
		naming[string(t)] = t
	}
	return EnumTable[T]{
		table[EnumRecord[T]]{name: tablename, PK: "id"},
		values,
		naming}
}

// Returns an array of each valid value of the composed enum type.
func (table EnumTable[T]) Values() []T {
	return table.values
}

func (table EnumTable[T]) Strings() []string {
	strings := make([]string, len(table.values))
	for i, enum := range table.values {
		strings[i] = string(enum)
	}
	return strings
}

func (table EnumTable[T]) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "name"  TEXT NOT NULL
  ) WITHOUT ROWID;`, table.Name)
}

func (table EnumTable[T]) SqlInit() string {
	rows := make([]string, len(table.values))
	for i, value := range table.values {
		rows[i] = fmt.Sprintf(`(%d, "%s")`, value, value)
	}
	return fmt.Sprintf(`INSERT INTO %s VALUES %s;`,
		table.Name, strings.Join(rows, ", "))
}

func (table EnumTable[T]) Insert(record *T) error {
	// TODO
	return nil
}

func (table EnumTable[T]) Get(id uint64) (EnumRecord[T], error) {
	// TODO
	return table.Zero(), nil
}

func (table EnumTable[T]) GetNamed(name string) (EnumRecord[T], error) {
	// TODO
	return table.Zero(), nil
}

func (table EnumTable[T]) SelectAll() (<-chan T, error) {
	// TODO
	return nil, nil
}

// The record abstraction for enums is interesting because
// it is both a single enum value (a one-byte int) and
// the unique (id, name) row in its backing database table.
type EnumRecord[T osn.EnumType] struct {
	value T
}

func (EnumRecord[T]) Columns() []string {
	return []string{"id", "name"}
}

func (record EnumRecord[T]) ScanRow(row *sql.Row) error {
	return errors.New("enums are immutable")
}

func (record EnumRecord[T]) ToValues() ([]driver.Value, error) {
	return []driver.Value{
		uint8(record.value), string(record.value),
	}, nil
}
