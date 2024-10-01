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
// github:kevindamm/wits-osn/map.go

package osn

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type LegacyMap struct {
	MapID int8   `json:"map_id"`
	Name  string `json:"name"`

	// The number of players this map can accommodate.
	// Use 0 for a player count on deprecated maps
	PlayerCount int `json:"player_count"`

	// Embedded type avoids the extra indirection
	// while facilitating compact table representation.
	LegacyMapDetails `json:"details,omitempty"`
}

func UnknownMap() LegacyMap {
	terrain := make([]byte, 0)
	units := make([]byte, 0)
	return LegacyMap{0, "UNKNOWN", 0,
		LegacyMapDetails{Filename: "", Terrain: terrain, Units: units}}
}

type LegacyMapDetails struct {
	Filename string `json:"-"`
	Theme    int    `json:"map_theme"`
	Terrain  []byte `json:"terrain,omitempty"`
	Units    []byte `json:"units,omitempty"`
	Width    int    `json:"columns"`
	Height   int    `json:"rows"`
}

// Hands the structure to a database driver using JSON serializagion.
func (details LegacyMapDetails) Value() (driver.Value, error) {
	return json.Marshal(details)
}

// Recovers the structure from a database driver using JSON deserialization.
func (details *LegacyMapDetails) Scan(value driver.Value) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid DB-value representation for LegacyMapDetails")
	}
	return json.Unmarshal(bytes, details)
}
