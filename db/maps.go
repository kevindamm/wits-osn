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

type LegacyMap osn.LegacyMap

func (LegacyMap) Columns() []string {
	return []string{"map_id", "map_name", "player_count", "details"}
}

func (record LegacyMap) ScanRow(row *sql.Row) error {
	bytes := make([]byte, 0)
	row.Scan(&record.MapID, &record.Name, &record.PlayerCount, &bytes)
	return json.Unmarshal(bytes, &record.LegacyMapDetails)
}

func (record LegacyMap) ToValues() ([]driver.Value, error) {
	return []driver.Value{
		record.MapID,
		record.Name,
		record.PlayerCount,
		record.LegacyMapDetails,
	}, nil
}

type tableMaps struct {
	table[LegacyMap]
	cached map[string]LegacyMap
}

func MakeMapsTable(sqldb *sql.DB) Table[LegacyMap] {
	return tableMaps{
		table: table[LegacyMap]{
			sqldb:   sqldb,
			name:    "maps",
			Primary: "map_id",
			NameCol: "map_name"}}
}

func (table tableMaps) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
      "map_id"        INTEGER PRIMARY KEY,
      "map_name"      VARCHAR(127) NOT NULL,
      "player_count"  INTEGER
			  CHECK(player_count == 2 OR player_count == 4 OR player_count == 0),

			-- Represents the following details in JSON-encoded bytes.
			"json"             BLOB
      -- "map_filename"  TEXT NOT NULL,
	    -- "theme"         INTEGER,
	    -- "terrain"       BLOB,  -- array of (coordinates and floor/wall)
	    -- "units"         BLOB,  -- array of (coordinates and unit type)
      -- "width"         INTEGER,
      -- "height"        INTEGER,

      FOREIGN KEY (map_theme) REFERENCES races (race_id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    ) WITHOUT ROWID;`, table.Name)
}

func (table tableMaps) SqlInit() string {
	return strings.Join([]string{
		fmt.Sprintf(`INSERT INTO %s VALUES
	    (0, "MAP_UNKNOWN", NULL, "", 0, NULL, 0, 0, true);`, table.Name),

		fmt.Sprintf(`INSERT INTO %s
      (map_id, map_name, player_count, map_filename, map_theme, width, height)
    VALUES
      (1,      "Machination",       4, "machination",        1,    13,    13),
      (2,      "Foundry (v1)",      0, "foundry",            1,    13,    12),
      (3,      "Foundry",           2, "foundry",            1,    13,    12),
      (4,      "Glitch",            2, "glitch",             1,    11,    11),
      (5,      "Candy Core Mine",   4, "candy-core-mine",    2,    13,    13),
      (6,      "Sweetie Plains",    2, "sweetie-plains",     2,    13,    13),
      (7,      "Peek-a-boo",        2, "peekaboo",           2,    13,    10),
      (8,      "Blitz Beach",       4, "blitz-beach",        3,    13,    11),
      (9,      "Long Nine",         2, "long-nine",          3,    13,    14),
      (10,     "Sharkfood Island",  2, "sharkfood-island",   3,    13,    10),
      (11,     "Acrospire",         4, "acrospire",          4,    13,    13),
      (12,     "Thorn Gulley",      2, "thorn-gulley",       4,    13,    12),
      (13,     "Reaper",            2, "reaper",             4,    13,    12),
      (14,     "Skull Duggery",     2, "skull-duggery",      3,    13,    10),
      (15,     "War Garden",        2, "war-garden",         4,    13,    12),
      (16,     "Sweet Tooth",       2, "sweet-tooth",        2,    13,    10),
      (17,     "Sugar Rock",        4, "sugar-rock",         2,    13,    13),
      (18,     "Mechanism",         4, "mechanism",          1,    13,    13)`,
			table.Name),
	}, ";\n")
}
