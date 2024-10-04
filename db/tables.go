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
	"fmt"
	"log"
	"strings"
)

// Abstraction over the structural type of values in a table or group of tables.
type Record interface {
	Columns() []string
	Values() ([]any, error)
	NamedValues() ([]driver.NamedValue, error)

	ScanValues(...driver.Value) error
	ScanRow(*sql.Row) error
	Scannables() []any
}

// The SQL-specific features of a relational table.  Fortunately, these have
// domains and codomains that aren't parameterized by the table's [Record] type,
// even though the computation of CREATE, INDEX and INSERT commands
type TableSql interface {
	Name() string
	SqlCreate() string
	SqlInit() string
}

// Common interface for one or more database tables and retrieval of records.
type Table[T Record] interface {
	TableSql
	Zero() T

	Get(int64) (T, error)
	GetByName(string) (T, error)
	SelectAll() (<-chan T, error)
}

// Variant on the table type which allows for inserting and deleting records.
type MutableTable[T Record] interface {
	Table[T]
	Insert(T) error
	Delete(int64) error
}

// Abstraction over one or more relational tables,
// represents an atomic structural type.
type tableBase[T Record] struct {
	sqldb *sql.DB

	// This table's name.
	name string

	// A zero instance of T, used for type reflection and instantiation.
	// It is assumed that nothing in this package exposes a pointer to it,
	// and it is semantically treated as immutable.
	zero T

	// It is a simplification for now, as these tables don't have complex queries.
	// There are a couple of places with joins but ScanRow() is still simple-ish.
	// Mainly the foreign key constraints and indices are a sanity check on data.

	// The primary key.  Defaults to `rowid` (Sqlite).
	Primary string

	// The column any "by name" lookups use.  Defaults to `name`
	NameCol string
}

func (table tableBase[T]) Name() string { return table.name }
func (table tableBase[T]) Zero() T      { return table.zero }

func (table tableBase[T]) Get(id int64) (T, error) {
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
	err = record.ScanRow(row)
	return record, err
}

// Retrieves the single row that corresponds to the indicated name.
func (table tableBase[T]) GetByName(name string) (T, error) {
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
	err = record.ScanRow(row)
	return record, err
}

func (table tableBase[T]) SelectAll() (<-chan T, error) {
	colnames := table.zero.Columns()
	query := fmt.Sprintf(`SELECT (%s) FROM %s;`,
		strings.Join(colnames, ", "), table.name)

	stmt, err := table.sqldb.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	channel := make(chan T)
	go func() {
		defer close(channel)
		defer rows.Close()

		for rows.Next() {
			var value T
			if rows.Err() != nil {
				return
			}
			rows.Scan(value.Scannables()...)
			channel <- value
		}
	}()

	return channel, nil
}

func qmarks(count int) string {
	str := make([]string, count)
	for i := range count {
		str[i] = "?"
	}
	return strings.Join(str, ", ")
}

type mutableBase[T Record] struct {
	tableBase[T]
}

// If the record's ID is 0 it will be updated with the index it was assigned.
// When a non-zero value is used as the record's primary key, if that record
// already existed it returns an error.
func (table mutableBase[T]) Insert(record T) error {
	colnames := record.Columns()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		table.name, colnames, qmarks(len(colnames)))
	stmt, err := table.sqldb.Prepare(query)
	if err != nil {
		return err
	}

	values, err := record.Values()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	fmt.Println()

	// TODO set the rowid value on `record`, if needed.
	// record.SetPrimary(exec_result.LastInsertId()) or similar
	return nil
}

func (table mutableBase[T]) Delete(id int64) error {
	primary_key := table.Primary
	if primary_key == "" {
		primary_key = "rowid"
	}
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;",
		table.name, primary_key)
	result, err := table.sqldb.Exec(sql, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows affected by DELETE(%s, %d)", table.name, id)
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
