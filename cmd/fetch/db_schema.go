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

package main

import (
	"errors"
	"log"
)

// Only needs to be called once at table setup
func CreateTables(db WitsDB) error {
	witsdb, ok := db.(*witsdb)
	if !ok {
		return errors.New("could not convert WitsDB to internal *witsdb")
	}
	for _, statement := range []string{
		// Enumerative relation between a map's name and its integer ID.
		`CREATE TABLE "maps" (
      "map_id"        INTEGER PRIMARY KEY,
      "map_name"      TEXT NOT NULL,
      "map_filename"  VARCHAR(127) NOT NULL,
      "player_count"  INTEGER CHECK(player_count == 2 OR player_count == 4),
      "deprecated"    BOOLEAN,

      UNIQUE (map_filename) ON CONFLICT IGNORE
    ) WITHOUT ROWID;`,

		`INSERT INTO maps 
      (map_id, map_name, map_filename, player_count, deprecated)
    VALUES
      (0,  "MAP_UNKNOWN", "", 2, true),
      (1,  "Machination", "machination", 4, false),
      (2,  "Foundry (v1)", "foundry-v1", 2, true),
      (3,  "Foundry", "foundry-v2", 2, false),
      (4,  "Glitch", "glitch", 2, false),
      (5,  "Candy Core Mine", "candy-core-mine", 4, false),
      (6,  "Sweetie Plains", "sweetie-plains", 2, false),
      (7,  "Peekaboo", "peekaboo", 2, false),
      (8,  "Blitz Beach", "blitz-beach", 4, false),
      (9,  "Long Nine", "long-nine", 2, false),
      (10, "Sharkfood Island", "sharkfood-island", 2, false),
      (11, "Acrospire", "acrospire", 4, false),
      (12, "Thorn Gulley", "thorn-gulley", 2, false),
      (13, "Reaper", "reaper", 2, false),
      (14, "Skull Duggery", "skull-duggery", 2, false),
      (15, "War Garden", "war-garden", 2, false),
      (16, "Sweet Tooth", "sweet-tooth", 2, false),
      (17, "Sugar Rock", "sugar-rock", 4, false),
      (18, "Mechanism", "mechanism", 4, false);`,

		// Enumerative relation for the different races.
		// Affects unit sprites and available hero unit.
		`CREATE TABLE "races" (
      "race_id" INTEGER PRIMARY KEY,
      "race_name" TEXT NOT NULL,

      UNIQUE (race_name) ON CONFLICT IGNORE
    ) WITHOUT ROWID;`,

		`INSERT INTO races (race_id, race_name) VALUES
      (0, "UNKNOWN"),
      (1, "Feedback"),
      (2, "Adorables"),
      (3, "Scallywags"),
      (4, "Veggienauts");`,

		// Enumerative relation for the different ranked leagues.
		// Players may be promoted or demoted based on win/loss records.
		`CREATE TABLE "leagues" (
      "league_id" INTEGER PRIMARY KEY,
      "league_name" TEXT NOT NULL,

      UNIQUE (league_name) ON CONFLICT IGNORE
    ) WITHOUT ROWID;`,

		`INSERT INTO leagues (league_id, league_name) VALUES
      (0, "UNKNOWN"),
      (1, "Fluffy"),
      (2, "Clever"),
      (3, "Gifted"),
      (4, "Master"),
      (5, "SuperTitan");`,

		// The player identifier and in-game display name,
		// used for foreign key relations in replay-related tables.
		`CREATE TABLE "players" (
      "id"    INTEGER PRIMARY KEY,
      "guid"  TEXT -- NOT NULL UNIQUE,
      "name"  TEXT NOT NULL
    ) WITHOUT ROWID;`,

		`INSERT INTO players (id, guid, name) VALUES
      (0, "", "UNKNOWN");`,

		// The `matches` metadata relates to an instance of a game between two
		// players.  This differs from the serialized replays that relate entirely
		// to each game's turns, or `solo_roles` which uniquely indexes the players
		// to their involvement in the match.
		`CREATE TABLE "matches" (
      -- rowid INTEGER PRIMARY KEY AUTOINCREMENT, -- legacy "id" or Index
      "match_hash"   TEXT NOT NULL,
      "competitive"  BOOLEAN,   -- league or friendly
      "season"       INTEGER,   -- seasons are of variable duration
      "start_time"   TIMESTAMP, -- time at creation, UTC

      "map_id"       INTEGER,   -- MapEnum
      "map_theme"    INTEGER,   -- matches RaceEnum
      "turn_count"   INTEGER,   -- number of turns (= one ply) for the match

      "version"      INTEGER,   -- the runtime version for this match
      "osn_status"   INTEGER,   -- this match's fetched -> parsed -> converted state

      FOREIGN KEY (map_id) REFERENCES maps (map_id)
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

      FOREIGN KEY (match_id) REFERENCES matches (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_id) REFERENCES players (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,

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

      "prev_points"   INTEGER DEFAULT 0,

      FOREIGN KEY (after) REFERENCES roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (until) REFERENCES roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_league) REFERENCES leagues (league_id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    );`,
	} {
		err := exec(witsdb.sqldb, statement)
		if err != nil {
			log.Println("error when executing SQL statement:")
			log.Println(statement)
			log.Fatal(err)
		}
	}
	return nil
}
