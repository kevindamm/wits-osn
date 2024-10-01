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
	"fmt"
	"log"
	"strings"
)

// Abstraction over the structural type of values in a table or group of tables.
type Record interface {
	Columns() []string

	ScanRow(*sql.Row) error
	ToValues() ([]driver.Value, error)
}

type TableSql interface {
	Name() string
	SqlCreate() string
	SqlInit() string
}

// Common interface for one or more database tables
type Table[T Record] interface {
	TableSql
	Zero() T

	Insert(*T) error
	Get(uint64) (T, error)
	GetByName(string) (T, error)
	SelectAll() *sql.Rows
}

// Only needs to be called once at database setup.  Also closes the database.
// Will LOG(FATAL) an error if creation or initialization fail, with SQL error.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	for _, table := range []TableSql{
		witsdb.status.table,
		witsdb.leagues.table,
		witsdb.races.table,
		witsdb.maps,
		witsdb.players,
		witsdb.matches,
		witsdb.roles,
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

// Abstraction over one or more relational tables,
// represents an atomic structural type.
type table[T Record] struct {
	sqldb *sql.DB
	name  string
	zero  T

	// Simplification for now, as these tables don't have complicated queries.
	// There are a couple of places with joins but ScanRow() is still simple-ish.
	// Mainly the foreign key constraints and indices are a sanity check on data.
	Primary string
	NameCol string
}

func (table table[T]) Name() string { return table.name }
func (table table[T]) Zero() T      { return table.zero }

func (table table[T]) Get(id uint64) (T, error) {
	var (
		record T
		err    error
	)
	rowid := table.Primary
	if rowid == "" {
		rowid = "rowid"
	}
	colstring := "(" + strings.Join(record.Columns(), ", ") + ")"

	sql := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = ?;`,
		colstring, table.name, rowid)

	row := table.sqldb.QueryRow(sql, id)
	if err != nil {
		return record, err
	}
	err = record.ScanRow(row)
	return record, err
}

// Retrieves the single row
func (table table[T]) GetByName(name string) (T, error) {
	var (
		record T
		err    error
	)
	namecol := table.NameCol
	if namecol == "" {
		namecol = "name"
	}
	colstring := "(" + strings.Join(record.Columns(), ", ") + ")"

	sql := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = ?;`,
		colstring, table.name, namecol)

	row := table.sqldb.QueryRow(sql, name)
	if err != nil {
		return record, err
	}
	err = record.ScanRow(row)
	return record, err
}

// If the record's ID is 0 it will be updated with the index it was assigned.
// When a non-zero value is used as the record's primary key, if that record
// already existed it returns an error.
func (table table[T]) Insert(record *T) error {
	// TODO
	return nil
}

// Returns the results of a `SELECT * FROM ?table.Name;` on the receiver ?table.
//
// If there is an error executing the "SELECT *" query (e.g., the table was
// dropped or not yet created) it will be surfaced in *sql.Rows which will
// always be non-nil (as per [database/sql.)
func (table table[T]) SelectAll() *sql.Rows {
	// TODO
	return nil
}

func (table table[T]) SqlCreate() string {
	// TODO derive this from Columns and Zero.NamedValues
	return fmt.Sprintf(`CREATE TABLE "%s" (
	  "%s" INTEGER PRIMARY KEY,
		"%s" TEXT NOT NULL
	) WITHOUT ROWID;`,
		table.Name, table.Primary, table.NameCol)
}

func (table table[T]) SqlInit() string {
	// TODO replace this method with explicit initialization in [OpenOsnDB] ctor.
	return ""
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
