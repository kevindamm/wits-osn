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
// github:kevindamm/wits-osn/db/maps.go

package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	osn "github.com/kevindamm/wits-osn"
)

type LegacyMapRecord osn.LegacyMap

func NewMapRecord() *LegacyMapRecord {
	record := new(LegacyMapRecord)
	record.LegacyMapDetails = osn.LegacyMapDetails{}
	return record
}

func (LegacyMapRecord) Columns() []string {
	return []string{"map_id", "map_name", "player_count", "filename"}
}

func (record LegacyMapRecord) Values() ([]any, error) {
	return []any{
		record.MapID,
		record.Name,
		record.PlayerCount,
		record.Filename,
	}, nil
}

func (record LegacyMapRecord) NamedValues() ([]driver.NamedValue, error) {
	return []driver.NamedValue{
		{
			Name:    "map_id",
			Ordinal: 0,
			Value:   record.MapID},
		{
			Name:    "map_name",
			Ordinal: 1,
			Value:   record.Name},
		{
			Name:    "player_count",
			Ordinal: 2,
			Value:   record.PlayerCount},
		{
			Name:    "details",
			Ordinal: 3,
			Value:   record.Filename},
	}, nil
}

func (record *LegacyMapRecord) ScanValues(values ...driver.Value) error {
	var ok bool
	record.MapID, ok = values[0].(uint8)
	if !ok {
		return fmt.Errorf("LegacyMap.MapID value %v not uint8", values[0])
	}
	record.Name, ok = values[1].(string)
	if !ok {
		return fmt.Errorf("LegacyMap.Name value %v not string", values[1])
	}
	record.PlayerCount, ok = values[2].(int)
	if !ok {
		return fmt.Errorf("LegacyMap.PlayerCount value %v not int", values[2])
	}
	bytes, ok := values[3].([]byte)
	if !ok {
		return fmt.Errorf("LegacyMap.Details value %v not []byte", values[3])
	}
	return json.Unmarshal(bytes, &record.LegacyMapDetails)
}

func (record *LegacyMapRecord) ScanRow(row *sql.Row) error {
	return row.Scan(&record.MapID, &record.Name, &record.PlayerCount, &record.Filename)
}

// This is provided for completeness but it has a flaw in that the details are
// not recoverable.  However, at present this method is only called from the
// `SELECT * FROM maps;` path, so it actually works out to get an abbreviation.
func (record *LegacyMapRecord) Scannables() []any {
	bytes := "{}"
	return []any{
		&record.MapID, &record.Name, &record.PlayerCount, &bytes,
	}
}

type tableMaps struct {
	tableBase[*LegacyMapRecord]
	cachedMaps map[string]osn.LegacyMap
}

func MakeMapsTable(sqldb *sql.DB) Table[*LegacyMapRecord] {
	return tableMaps{
		tableBase: tableBase[*LegacyMapRecord]{
			sqldb:   sqldb,
			name:    "maps",
			zero:    NewMapRecord(),
			new:     NewMapRecord,
			Primary: "map_id",
			NameCol: "map_name"},
		cachedMaps: make(map[string]osn.LegacyMap)}
}

func (table tableMaps) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
      "map_id"        INTEGER PRIMARY KEY,
      "map_name"      VARCHAR(127) NOT NULL,
      "player_count"  INTEGER
			  CHECK(player_count == 2 OR player_count == 4 OR player_count == 0),

			-- Represents the following details in JSON-encoded bytes.
			"filename"      TEXT -- see LegacyMapDetails
    ) WITHOUT ROWID;`, table.name)
}

func (table tableMaps) SqlInit() string {
	return strings.Join([]string{
		fmt.Sprintf(`CREATE UNIQUE INDEX map_names ON %s (map_name);`,
			table.name),

		fmt.Sprintf(`INSERT INTO %s VALUES
	    (0, "MAP_UNKNOWN", NULL, "")`, table.name),

		fmt.Sprintf(`INSERT INTO %s
      (map_id, map_name, player_count, filename)
    VALUES
      (1,      "Machination",       4, "machination.json"),
      (2,      "Foundry (v1)",      0, ""),
      (3,      "Foundry",           2, "foundry.json"),
			(4,      "Glitch",            2, "glitch.json"),
			(5,      "Candy Core Mine",   4, "candy-core-mine.json"),
			(6,      "Sweetie Plains",    2, "sweetie-plains.json"),
			(7,      "Peek-a-boo",        2, "peekaboo.json"),
			(8,      "Blitz Beach",       4, "blitz-beach.json"),
			(9,      "Long Nine",         2, "long-nine.json"),
			(10,     "Sharkfood Island",  2, "sharkfood-island.json"),
			(11,     "Acrospire",         4, "acrospire.json"),
			(12,     "Thorn Gulley",      2, "thorn-gulley.json"),
			(13,     "Reaper",            2, "reaper.json"),
			(14,     "Skull Duggery",     2, "skull-duggery.json"),
			(15,     "War Garden",        2, "war-garden.json"),
			(16,     "Sweet Tooth",       2, "sweet-tooth.json"),
			(17,     "Sugar Rock",        4, "sugar-rock.json"),
			(18,     "Mechanism",         4, "mechanism.json")`,
			table.name),
	}, ";\n")
}
