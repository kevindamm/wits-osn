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
	SqlCreate(string) string
	SqlInit(string) string

	IsValid() bool
	ScanRecord(*sql.Row) error
	RecordValues() ([]driver.Value, error)
}

// Abstraction over one or more tables, represents an atomic structural type.
type Table[T Record] struct {
	Zero    T
	Name    string
	Primary string
}

type TableIndex[T Record, Index comparable] struct {
	Table      *Table[T]
	Column     string
	ColumnZero Index
}

func (table Table[T]) CreateAndInit(db *sql.DB) {
	log.Println("creating table '" + table.Name + "'")
	must_execsql(db, table.Zero.SqlCreate(table.Name))

	initsql := table.Zero.SqlInit(table.Name)
	if initsql != "" {
		log.Println("populating table '" + table.Name + "'")
		must_execsql(db, initsql)
	}
}

// Only needs to be called once at database setup.  Also closes the database.
// Will LOG(FATAL) an error if creation or initialization fail, with SQL error.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	witsdb.status.CreateAndInit(witsdb.sqldb)

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

// Execute the SQL statement on the provided database.
//
// Prefer the public methods of the WitsDB interface, this is intentionally
// reserved for one-off statements where the result-set isn't important.
func execsql(db *sql.DB, sql string) error {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

// Convenience check-fail for errors encountered during database setup.
// Fails loudly so that awareness of likely schema corruption is foremost.
// Reserved for package-internal because it should only be used during setup.
func must_execsql(db *sql.DB, sql string) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println("error preparing SQL")
		log.Println(sql)
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println("error executing SQL")
		log.Println(sql)
		log.Fatal(err)
	}
}
