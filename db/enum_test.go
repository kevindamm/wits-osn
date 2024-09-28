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

	"github.com/kevindamm/wits-osn/db"
)

type MyEnumType uint8

const (
	ENUM_UNKNOWN MyEnumType = iota
	ENUM_AAA
	ENUM_BEE
	ENUM_SEE
	EnumRange
	ENUM_KAY = uint8('K') // TODO test noncontiguous range
)

var enum_names = []string{
	"UNKNOWN",
	/* AAA */ "Hey!",
	/* BEE */ "Are, 'e a tee ache eee",
	/* SEE */ "Sí?",
	// /* KAY */ "ok",
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
	yay := ENUM_AAA

	// The same value can be interpreted as a number or via String()-ification.
	formatted := fmt.Sprintf("%d: %s", yay, yay)

	if formatted != "1: Hey!" {
		t.Errorf("improper casting, formatted = '%s'", formatted)
	}
}

func TestDBEnumTable(t *testing.T) {
	tablename := "myenum"
	enumtable := db.MakeEnumTable(tablename, EnumRange)
	create := enumtable.Zero.SqlCreate(tablename)
	if create != fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "name"  VARCHAR(9) NOT NULL
  ) WITHOUT ROWID;`, tablename) {
		t.Errorf("SQL CREATE TABLE malformed; got\n%s", create)
	}

	init := enumtable.Zero.SqlInit(tablename)
	if create != fmt.Sprintf(`INSERT INTO %s VALUES (0, "%s"), (1, "%s"), (2, "%s"), (3, "%s");`,
		tablename, enum_names[0], enum_names[1], enum_names[2], enum_names[3]) {
		t.Errorf("SQL INSERT VALUES malformed; got\n%s", init)
	}
}

type GenericEnum interface {
	~uint8
	String() string
}

type Record interface {
	Columns(string) string
}

type EnumRecord[EnumType GenericEnum] struct {
	Value EnumType
	Range EnumType
}

func (enum EnumRecord[T]) Columns(string) string {
	return fmt.Sprintf("columns []string\n  %s to %s\n",
		enum.Value.String(), T(enum.Range-1))
}

type EnumTable[EnumType GenericEnum] struct {
	Name string
	Zero EnumRecord[EnumType]
}

func MakeEnumTable[EnumType GenericEnum](tablename string, enumrange EnumType) EnumTable[EnumType] {
	return EnumTable[EnumType]{
		Name: tablename,
		Zero: EnumRecord[EnumType]{Range: enumrange},
	}
}

func TestEnumTable(t *testing.T) {
	tablename := "myenum"
	enumtable := MakeEnumTable(tablename, EnumRange)
	//expected := fmt.Sprintf("INSERT INTO %s VALUES (0, %s), (1, %s), (2, %s), (3, %s);",
	//	tablename, enum_names[0], enum_names[1], enum_names[2], enum_names[3]) //, enum_names[4])
	expected := `columns []string
  UNKNOWN to Sí?
`

	if enumtable.Zero.Columns(tablename) != expected {
		t.Errorf("init SQL doesn't match expectation\n%s\n%s",
			enumtable.Zero.Columns(tablename), expected)
	}
}
