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

	osn "github.com/kevindamm/wits-osn"
)

// Abstraction over the structural type of values in a table or group of tables.
type Record interface {
	SqlCreate(string) string
	SqlInit(string) []string

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

// Serialized enum types are a recurring pattern in fully defined relational DBs
// and here we make the simplifying assumption that any enumerated type can fit
// inside a single byte.  If the domain is larger than 255 items, express that
// relation as a [Table[T Record]], named by a [TableIndex[Record, comparable]].
type EnumTable[T osn.EnumType] struct {
	Name string

	// Represents the limit value (excl.) for the valid enumeration range.
	// Enum definitions should use const, iota and a trailing enum value to
	// derive this value from the last valid enumeration.
	limit T
}

// Constructor for EnumTable[T] that can be called from other modules.
// However, it is only ever called from the database initialization or tests.
//
// EnumTable is a common pattern of (id, name) pairwise relations.
// The `id` is its enumeration and the `name` is persisted in the DB for
// documentation.  The list of pairs is provided by a Functor: int -> string.
func MakeEnumTable[T osn.EnumType](tablename string, enumrange T) EnumTable[T] {
	return EnumTable[T]{tablename, enumrange}
}

// Returns an array of each valid value of the composed enum type.
func (table EnumTable[T]) Values() []T {
	values := make([]T, table.limit)
	for i := range table.limit {
		values[i] = T(i)
	}
	return values
}

func (table EnumTable[T]) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "name"  TEXT NOT NULL
  ) WITHOUT ROWID;`, table.Name)
}

func (table EnumTable[T]) SqlInit() string {
	rows := make([]string, table.limit)
	for i, value := range table.Values() {
		rows[i] = fmt.Sprintf(`(%d, "%s")`, value, value)
	}
	return fmt.Sprintf(`INSERT INTO %s VALUES %s;`,
		table.Name, strings.Join(rows, ", "))
}

type TableIndex[T Record, Index comparable] struct {
	Table      *Table[T]
	Column     string
	ColumnZero Index
}

func init_enum_table[T osn.EnumType](db *sql.DB, table EnumTable[T]) {
	log.Println("creating table '" + table.Name + "'")
	must_execsql(db, table.SqlCreate())

	initsql := table.SqlInit()
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

	init_enum_table(witsdb.sqldb, witsdb.status)
	init_enum_table(witsdb.sqldb, witsdb.leagues)
	init_enum_table(witsdb.sqldb, witsdb.races)

	// This could be improved; the table name is included in table instantiation.
	log.Print(witsdb.maps.Zero.SqlCreate("maps"))
	must_execsql(witsdb.sqldb, witsdb.maps.Zero.SqlCreate("maps"))

	for _, sql := range witsdb.maps.Zero.SqlInit("maps") {
		log.Print(sql)
		must_execsql(witsdb.sqldb, sql)
	}

	for _, schema := range [][]string{
		// map metadata, map data contained in JSON files
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
