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
	"strings"
)

// Abstraction over the structural type of values in a table or group of tables.
type Record interface {
	Columns() []string

	ScanRow(*sql.Row) error
	ToValues() ([]driver.Value, error)
}

// Common interface for one or more database tables
type Table[T Record] interface {
	Name() string
	Zero() Record

	SqlCreate() string
	SqlInit() string

	Insert(*T) error
	Get(id uint64) (T, error)
	GetNamed(name string) (T, error)
	SelectAll() (<-chan T, error)
}

// Abstraction over one or more relational tables,
// represents an atomic structural type.
type table[T Record] struct {
	name string
	zero T
	PK   string
}

func (table table[T]) Name() string { return table.name }
func (table table[T]) Zero() T      { return table.zero }

// Only needs to be called once at database setup.  Also closes the database.
// Will LOG(FATAL) an error if creation or initialization fail, with SQL error.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	for _, table := range []Table[Record]{
		witsdb.status,
		witsdb.leagues,
		witsdb.races,
		witsdb.maps,
		witsdb.players,
		witsdb.matches,
		witsdb.player_roles,
		witsdb.standings,
	} {
		log.Println("creating table '" + table.Name() + "'")
		must_execsql(witsdb.sqldb, table.SqlCreate())

		initsql := table.SqlInit()
		if initsql != "" {
			log.Println("populating table '" + table.Name() + "'")
			statements := strings.Split(initsql, ";")
			for _, sql := range statements {
				must_execsql(witsdb.sqldb, sql)
			}
		}
	}

	return nil
}

// Execute the SQL statement on the provided database.
//
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
