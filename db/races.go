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
// github:kevindamm/wits-osn/db/races.go

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Enumeration of race; determines special unit and affects visual appearance.
// Satisfies Resource[RaceEnum] for inclusion in database tables.
type RaceEnum uint8

func RaceEnumName(race RaceEnum) string {
	if int(race) >= len(race_names) {
		return race_names[RACE_UNKNOWN]
	}
	return race_names[race]
}

func RaceEnumSpecial(race RaceEnum) string {
	if int(race) >= len(race_names) {
		return race_specials[RACE_UNKNOWN]
	}
	return race_specials[race]
}

const (
	RACE_UNKNOWN RaceEnum = iota
	RACE_FEEDBACK
	RACE_ADORABLES
	RACE_SCALLYWAGS
	RACE_VEGGIENAUTS
)

var race_names = [...]string{
	"",
	"Feedback",
	"Adorables",
	"Scallywags",
	"Veggienauts",
}

var race_specials = map[RaceEnum]string{
	RACE_UNKNOWN:     "Unknown",
	RACE_FEEDBACK:    "Scrambler",
	RACE_ADORABLES:   "Mobi",
	RACE_SCALLYWAGS:  "Bombshell",
	RACE_VEGGIENAUTS: "Bramble",
}

func (race RaceEnum) race_values() []string {
	values := make([]string, len(race_names))
	for i, name := range race_names {
		values[i] = fmt.Sprintf(`(%2d,  "%s",  "%s")`,
			i, name, RaceEnumSpecial(RaceEnum(i)))
	}
	return values
}

func (race RaceEnum) Name() string { return "races" }
func (race RaceEnum) CreateAndInitSql(db *sql.DB) error {
	for _, sql := range []string{
		`create TABLE "races" (
      "race_id"    INTEGER PRIMARY KEY,
      "race_name"  VARCHAR(12) UNIQUE NOT NULL,
      "special"    VARCHAR(10) UNIQUE NOT NULL
    ) WITHOUT ROWID;`,

		"INSERT INTO races VALUES\n  " + strings.Join(
			RaceEnum(0).race_values(), ",\n  ") + ";",
	} {
		err := execsql(db, sql)
		if err != nil {
			log.Println("error in [CreateTablesFromResource] when executing SQL:")
			log.Println(sql)
			return err
		}
	}
	return nil
}

func (race RaceEnum) SqlGet(id int64) string {
	return fmt.Sprintf("SELECT * FROM maps WHERE id = %d", id)
}

func (race RaceEnum) ScanRecord(*sql.Row) error {
	return errors.New(
		"attempted to insert a new RaceEnum value; table is readonly")
}
