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
// github:kevindamm/wits-osn/cmd/fetch/db_schema.go

package db

import (
	"database/sql"
	"errors"
	"log"
)

// Only needs to be called once at table setup.  Also closes the database.
func CreateTablesAndClose(db_path string) error {
	witsdb, err := open_database(db_path)
	if err != nil {
		return errors.New("could not open WitsDB database")
	}
	defer witsdb.Close()

	for _, statement := range []string{
		// Enumerative relation for unit races.
		`create TABLE "races" (
      "race_id"    INTEGER PRIMARY KEY,
      "race_name"  VARCHAR(20) UNIQUE NOT NULL,
      "hero"       VARCHAR(10) UNIQUE NOT NULL
    ) WITHOUT ROWID;`,

		`INSERT INTO races VALUES
      (1, "Feedback",    "Scrambler"),
      (2, "Adorables",   "Mobi"),
      (3, "Scallywags",  "Bombshell"),
      (4, "Veggienauts", "Bramble");`,

		// Enumerative relation for maps, including some additional properties.
		`CREATE TABLE "maps" (
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
    ) WITHOUT ROWID;`,

		`INSERT INTO maps VALUES (0, "MAP_UNKNOWN", NULL, "", 0, NULL, 0, 0, true);`,
		`INSERT INTO maps 
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
		`UPDATE maps set deprecated = true WHERE map_id = 2;`,

		// Enumerative relation for the different ranked leagues.
		// Players may be promoted or demoted based on win/loss records.
		`CREATE TABLE "leagues" (
      "league_id"    INTEGER PRIMARY KEY,
      "league_name"  TEXT NOT NULL,

      UNIQUE (league_name) ON CONFLICT IGNORE
    ) WITHOUT ROWID;`,

		`INSERT INTO leagues VALUES
      (0, "UNKNOWN"),
      (1, "Fluffy"),
      (2, "Clever"),
      (3, "Gifted"),
      (4, "Master"),
      (5, "SuperTitan");`,

		`CREATE TABLE "fetch_status" (
      "id"    INTEGER PRIMARY KEY,
      "name"  VARCHAR(10) NOT NULL
    ) WITHOUT ROWID;`,

		`INSERT INTO fetch_status VALUES
      (0, "UNKNOWN"),
      (1, "LISTED"),
      (2, "FETCHED"),
      (3, "UNWRAPPED"),
      (4, "CONVERTED"),
      (5, "CANONICAL"),
      (6, "VALIDATED"),
      (7, "INDEXED"),
      (8, "INVALID"),
      (9, "LEGACY");`,

		// The player identifier and in-game display name,
		// used for foreign key relations in replay-related tables.
		`CREATE TABLE "players" (
      "id"    INTEGER PRIMARY KEY,
      "name"  TEXT NOT NULL
    ) WITHOUT ROWID;`,

		`INSERT INTO players (id, name) VALUES (0, "UNKNOWN");`,

		// An index is created to assist lookup of players by name.
		`CREATE UNIQUE INDEX player_names
      ON players (name);`,

		// We don't always know the GC:ID (it is revealed in replays),
		// having it as a separate table affords
		// both optional and non-null uniqueness.
		`CREATE TABLE "player_gcid" (
      "player_id"  INTEGER PRIMARY KEY,
      "gcid"       TEXT NOT NULL UNIQUE,

      FOREIGN KEY (player_id)
        REFERENCES players (id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    ) WITHOUT ROWID;`,

		// The `matches` metadata relates to an instance of a game between two
		// players.  This differs from the serialized replays that relate entirely
		// to each game's turns, or `solo_roles` which uniquely indexes the players
		// to their involvement in the match.
		`CREATE TABLE "matches" (
      -- rowid INTEGER PRIMARY KEY AUTOINCREMENT, -- legacy "id" or Index
      "match_hash"   TEXT NOT NULL UNIQUE,
      "match_index"  INTEGER,      -- OSN row id
      "competitive"  BOOLEAN,      -- league or friendly
      "season"       INTEGER,      -- seasons are of variable duration
      "start_time"   TIMESTAMP,    -- time at creation, UTC

      "map_id"       INTEGER,      -- MapEnum
      "turn_count"   INTEGER,      -- number of turns (= one ply) for the match

      "version"      INTEGER,      -- the runtime version for this match
      "status"       VARCHAR(10),  -- this match's fetch status

      FOREIGN KEY (map_id)
        REFERENCES maps (map_id)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (status)
        REFERENCES fetch_status (name)
        ON DELETE CASCADE ON UPDATE NO ACTION
    );`,

		// Relation for which players are participating in which matches,
		// and the turn order they are assigned to.
		// Appropriate for both 1v1 and 2v2 matches.
		`CREATE TABLE "roles" (
      -- rowid INTEGER PRIMARY KEY,
      "match_id" INTEGER NOT NULL,
      "player_id" INTEGER NOT NULL,
      "turn_order" INTEGER CHECK(turn_order > 0 AND turn_order <= 2),

      FOREIGN KEY (match_id)
        REFERENCES matches (rowid)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
      FOREIGN KEY (player_id)
        REFERENCES players (rowid)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,

      UNIQUE (match_id, turn_order) ON CONFLICT FAIL,
      UNIQUE (match_id, player_id) ON CONFLICT IGNORE
    );`,

		// Represents the ranked standings for a player at a point in time.
		`CREATE TABLE "standings" (
      -- rowid INTEGER PRIMARY KEY,
      "after" INTEGER NOT NULL UNIQUE,
      "until" INTEGER,

      "player_league" INTEGER NOT NULL,
      "player_rank"   INTEGER NOT NULL,
      "player_points" INTEGER DEFAULT 0,
      "player_delta"  INTEGER DEFAULT 0,

      FOREIGN KEY (after)
        REFERENCES roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (until)
        REFERENCES roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_league)
        REFERENCES leagues (league_id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    );`,
	} {
		err := execsql(witsdb.sqldb, statement)
		if err != nil {
			log.Println("error when executing SQL statement:")
			log.Println(statement)
			log.Fatal(err)
		}
	}
	return nil
}

// Convenience function for ad-hoc query execution.
func execsql(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}
