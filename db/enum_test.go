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
// github:kevindamm/wits-osn/db/enum_test.go

package db_test

import (
	"fmt"
	"testing"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

type MyEnumType uint8

const (
	ENUM_UNKNOWN MyEnumType = iota
	ENUM_AYY
	ENUM_BEE
	ENUM_SEE
	EnumRange
)

var enum_names = []string{
	"UNKNOWN",
	/* AYY */ "Hey!",
	/* BEE */ "Be, not to not be.",
	/* SEE */ "Sí?",
}

func (enum MyEnumType) IsValid() bool {
	return uint8(enum) < uint8(EnumRange)
}

func (enum MyEnumType) String() string {
	if !enum.IsValid() {
		return ENUM_UNKNOWN.String()
	}
	return enum_names[enum]
}

// String conversion uses the enum's type (MyEnumType) not the underlying uint8.
func TestEnum(t *testing.T) {
	yay := ENUM_AYY

	// The same value can be interpreted as a number or via String()-ification.
	formatted := fmt.Sprintf("%d: %s", yay, yay)

	if formatted != "1: Hey!" {
		t.Errorf("improper casting, %d formatted = '%s'",
			yay, formatted)
	}
}

func TestDBEnumTable(t *testing.T) {
	tablename := "myenum"
	enumtable := db.MakeEnumTable(tablename, osn.EnumValuesFor(EnumRange))
	create := enumtable.SqlCreate()
	expected := `CREATE TABLE "myenum" (
     "id"    INTEGER PRIMARY KEY,
     "name"  TEXT NOT NULL
   ) WITHOUT ROWID;`
	if create != expected {
		t.Errorf("SQL CREATE TABLE malformed; got\n%s\nexpected\n%s\n", create, expected)
	}

	init := enumtable.SqlInit()
	expected = fmt.Sprintf(`INSERT INTO myenum VALUES (0, "%s"), (1, "%s"), (2, "%s"), (3, "%s");`,
		enum_names[0], enum_names[1], enum_names[2], enum_names[3])
	if init != expected {
		t.Errorf("SQL INSERT VALUES malformed; got\n%s\nwant\n%s", init, expected)
	}
}
