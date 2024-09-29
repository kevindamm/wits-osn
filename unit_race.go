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
// github:kevindamm/wits-osn/unit_race.go

package osn

// Enumeration of race; determines special unit and affects visual appearance.
// Satisfies Resource[RaceEnum] for inclusion in database tables.
type UnitRaceEnum uint8

const (
	RACE_UNKNOWN UnitRaceEnum = iota
	RACE_FEEDBACK
	RACE_ADORABLES
	RACE_SCALLYWAGS
	RACE_VEGGIENAUTS
	UnitRaceRange
)

func (race UnitRaceEnum) IsValid() bool {
	return race < UnitRaceRange
}

var race_names = []string{
	"UNITRACE_UNKNOWN",
	"Feedback",
	"Adorables",
	"Scallywags",
	"Veggienauts",
}

func (race UnitRaceEnum) String() string {
	if race.IsValid() {
		return race_names[race]
	}
	return RACE_UNKNOWN.String()
}

// Each race has one special unit associated with it.
// They share the enumeration ordering with UnitRaceEnum for ease of conversion.
type UnitSpecialEnum uint8

const (
	SPECIAL_UNKNOWN UnitSpecialEnum = iota
	SPECIAL_SCRAMBLER
	SPECIAL_MOBI
	SPECIAL_BOMBSHELL
	SPECIAL_BRAMBLE
	UnitSpecialRange
)

func (special UnitSpecialEnum) IsValid() bool {
	return special < UnitSpecialRange
}

var special_names = []string{
	"Unknown",
	"Scrambler",
	"Mobi",
	"Bombshell",
	"Bramble",
}

func (special UnitSpecialEnum) String() string {
	if special.IsValid() {
		return special_names[special]
	}
	return SPECIAL_UNKNOWN.String()
}
