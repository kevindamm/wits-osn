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
	"log"

	osn "github.com/kevindamm/wits-osn"
)

type LegacyMap osn.LegacyMap

func (LegacyMap) SqlCreate(tablename string) string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
      "map_id"        INTEGER PRIMARY KEY,
      "map_name"      VARCHAR(127) NOT NULL,
      "player_count"  INTEGER CHECK(player_count == 2 OR player_count == 4),

      -- The following could be collected into a JSON field.
      "map_filename"  TEXT NOT NULL,
      "map_theme"     INTEGER,
      "map_json"      TEXT,
      "width"         INTEGER,
      "height"        INTEGER,

      "deprecated"    BOOLEAN DEFAULT FALSE,

      FOREIGN KEY (map_theme) REFERENCES races (race_id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    ) WITHOUT ROWID;`, tablename)
}

func (LegacyMap) SqlInit(tablename string) []string {
	return []string{
		fmt.Sprintf(`INSERT INTO %s VALUES
	    (0, "MAP_UNKNOWN", NULL, "", 0, NULL, 0, 0, true);`, tablename),

		fmt.Sprintf(`INSERT INTO %s
      (map_id, map_name, player_count, map_filename, map_theme, width, height)
    VALUES
      (1,      "Machination",       4, "machination",        1,    13,    13),
      (2,      "Foundry (v1)",      2, "foundry",            1,    13,    12),
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
      (18,     "Mechanism",         4, "mechanism",          1,    13,    13);`,
			tablename),

		fmt.Sprintf(`UPDATE %s set deprecated = true WHERE map_id = 2;`,
			tablename),
	}
}

func (gamemap LegacyMap) Details() []byte {
	bytes, err := json.Marshal(gamemap.LegacyMapDetails)
	if err != nil {
		log.Fatalf("error marshaling legacy map details\n%s", err)
	}
	return bytes
}

// TODO
func (LegacyMap) IsValid() bool { return true }

func (gamemap LegacyMap) ScanRecord(*sql.Row) error {
	return nil
}

func (gamemap LegacyMap) RecordValues() ([]driver.Value, error) {
	return []driver.Value{
		gamemap.MapID, gamemap.Name, gamemap.PlayerCount, gamemap.Details(),
	}, nil
}

type MapDetails struct {
	Filename string           `json:"-"`
	Theme    osn.UnitRaceEnum `json:"map_theme"`
	Terrain  []byte           `json:"terrain,omitempty"`
	Units    []byte           `json:"units,omitempty"`
	Width    int              `json:"columns"`
	Height   int              `json:"rows"`
}

func (db *osndb) AllMaps() ([]osn.LegacyMap, error) {
	maps := make([]osn.LegacyMap, 0)

	stmt, err := db.sqldb.Prepare(`SELECT
	  map_id, map_name, player_count, map_filename, map_theme, width, height
	  FROM maps
		WHERE NOT deprecated;`)
	if err != nil {
		return maps, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return maps, err
	}
	defer rows.Close()

	mapobj := osn.LegacyMap{}
	for rows.Next() {
		// map_id, map_name, player_count, map_filename, map_theme, width, height
		rows.Scan(
			&mapobj.MapID, &mapobj.Name, &mapobj.PlayerCount,
			&mapobj.Filename, &mapobj.Theme,
			&mapobj.Width, &mapobj.Height)
		maps = append(maps, mapobj)
	}

	return maps, nil
}
