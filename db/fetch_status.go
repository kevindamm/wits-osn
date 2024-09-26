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
// github:kevindamm/wits-osn/db/fetch_status.go

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	osn "github.com/kevindamm/wits-osn"
)

type FetchStatus osn.FetchStatus

// Enumerative relation for the different steps of the fetcher/crawler.
func (FetchStatus) SqlCreate(tablename string) string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "name"  VARCHAR(9) NOT NULL
  ) WITHOUT ROWID;`, tablename)
}

// Populates the table with all FetchStatus enum values.
func (FetchStatus) SqlInit(tablename string) string {
	values := make([]string, osn.FetchStatusRange)
	for i := range osn.FetchStatusRange {
		status := osn.FetchStatus(i)
		values[i] = fmt.Sprintf(`(%d, "%s")`, status, status)
	}
	return fmt.Sprintf(`INSERT INTO %s VALUES %s;`,
		tablename, strings.Join(values, ", "))
}

func (record FetchStatus) IsValid() bool {
	return uint(record) < uint(osn.FetchStatusRange)
}

// This is intentionally a value-receiver instead of a pointer-receiver.
func (record FetchStatus) ScanRecord(row *sql.Row) error {
	return row.Scan(record)
}

func (record FetchStatus) RecordValues() ([]driver.Value, error) {
	if osn.FetchStatus(record) >= osn.FetchStatusRange {
		return FetchStatus(osn.STATUS_UNKNOWN).RecordValues()
	}
	return []driver.Value{
		int(record), osn.FetchStatus(record).String(),
	}, nil
}

func FetchStatusByName(table *Table[FetchStatus]) TableIndex[FetchStatus, string] {
	return TableIndex[FetchStatus, string]{Table: table, Column: "name"}
}
