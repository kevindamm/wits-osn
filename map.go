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

type LegacyMap struct {
	MapID uint8  `json:"map_id"`
	Name  string `json:"name"`

	// The number of players this map can accommodate.
	// Use 0 for a player count on deprecated maps
	PlayerCount int `json:"player_count"`

	// Embedded type avoids the extra indirection
	// while facilitating compact table representation.
	LegacyMapDetails
}

func UnknownMap() LegacyMap {
	return LegacyMap{0, "UNKNOWN", 0,
		LegacyMapDetails{Filename: ""}}
}

type LegacyMapDetails struct {
	Filename string         `json:"filename"`
	Width    int            `json:"columns"`
	Height   int            `json:"rows"`
	Theme    int            `json:"theme"`
	Defaults map[string]int `json:"defaults"`
	Init     []MapTileInit  `json:"background"`
}

type MapTileInit struct {
	I     int         `json:"i"`
	J     int         `json:"j"`
	Type  MapTileType `json:"type"`
	Sub   SpriteIndex `json:"sub,omitempty"`
	Owner PlayerIndex `json:"owner,omitempty"`
}

type MapTileType uint8
type SpriteIndex uint8
type PlayerIndex uint8

//// Hands the structure to a database driver using JSON serializagion.
//func (details LegacyMapDetails) Value() (driver.Value, error) {
//	return json.Marshal(details)
//}

//// Recovers the structure from a database driver using JSON deserialization.
//func (details *LegacyMapDetails) Scan(value driver.Value) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New("invalid DB-value representation for LegacyMapDetails")
//	}
//	return json.Unmarshal(bytes, details)
//}
