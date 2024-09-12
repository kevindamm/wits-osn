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
// github:kevindamm/wits-osn/units.go

package osn

import "encoding/json"

// OSN enumeration of the different squads a player can choose from.
type UnitRaceEnum byte

const (
	RACE_UNKNOWN UnitRaceEnum = iota
	RACE_FEEDBACK
	RACE_ADORABLES
	RACE_SCALLYWAGS
	RACE_VEGGIENAUTS
)

func (race UnitRaceEnum) String() string {
	if race > RACE_VEGGIENAUTS {
		race = RACE_UNKNOWN
	}
	return []string{
		"UNKNOWN",
		"FEEDBACK",
		"ADORABLES",
		"SCALLYWAGS",
		"VEGGIENAUTS",
	}[int(race)]
}

func (race UnitRaceEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(race))
}
