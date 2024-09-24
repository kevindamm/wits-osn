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
	"errors"
	"log"
)

type Table[T Resource] struct {
	Name string
	Type T
}

func Columns() []string {

	return []string{"rowid", ""}
}

// Abstraction over one or more tables, representing an atomic structural type.
//
// Some are defined in `base.go` (Maps, Leagues, Races, etc.) while others are
// defined in the files associated with their type (Player, Match, Replay...).
type Resource interface {
	Name() string
	CreateAndInitSql(db *sql.DB) error
	SqlGet(db *sql.DB, id int64) error

	ScanRecord(*sql.Row) error
	Columns() []string
	//Values() []sql.ColumnType
	//Subset([]string) []sql.ColumnType
}

func ExecBatch(db *osndb, schemata ...[]string) error {
	for _, schema := range schemata {
		for _, sql := range schema {
			err := execsql(db.sqldb, sql)
			if err != nil {
				log.Println("error when executing SQL statement:")
				log.Println(sql)
				return err
			}
		}
	}
	return nil
}

// Only needs to be called once at table setup.  Also closes the database.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	err = ExecBatch(witsdb, [][]string{
		FetchStatusSchema, LeaguesSchema, RacesSchema,
		LegacyMapSchema,
		PlayerSchema,
		MatchesSchema,
		PlayerRoleSchema,
		PlayerStandingsSchema,
	}...)
	return err
}

// Convenience function for ad-hoc query execution.
func execsql(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
