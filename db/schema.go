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
// github:kevindamm/wits-osn/db/schema.go

package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
)

// Abstraction over the structural type of values in a table or group of tables.
type Record interface {
	SqlCreateAndInit() string

	ScanRecord(*sql.Row) error
	RecordValues() ([]driver.Value, error)
}

// Abstraction over one or more tables, represents an atomic structural type.
type Table[T Record] struct {
	Zero    T
	Name    string
	Primary string
}

// Only needs to be called once at database setup.  Also closes the database.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	err = execsql(witsdb.sqldb, witsdb.status.Zero.SqlCreateAndInit())
	if err != nil {
		return err
	}

	for _, schema := range [][]string{
		// base enums
		LeaguesSchema, RacesSchema,
		// map metadata, map data contained in JSON files
		LegacyMapSchema,
		PlayerSchema,
		// matches must be defined before roles and standings
		MatchesSchema,
		PlayerRoleSchema,
		PlayerStandingsSchema,
	} {
		for _, sql := range schema {
			err := execsql(witsdb.sqldb, sql)
			if err != nil {
				log.Println("when creating and initializing tables, SQL:")
				log.Println(sql)
				return err
			}
		}
	}
	return nil
}

func execsql(db *sql.DB, sql string) error {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
