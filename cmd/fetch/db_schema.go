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
      "deprecated"    BOOLEAN,

      UNIQUE (map_filename) ON CONFLICT IGNORE
    ) WITHOUT ROWID;`,

		`INSERT INTO maps 
      (map_id, map_name, map_filename, deprecated)
    VALUES
      (0,  "MAP_UNKNOWN", "", true),
      (1,  "Machination", "machination", false),
      (2,  "Foundry (v1)", "foundry-v1", true),
      (3,  "Foundry", "foundry-v2", false),
      (4,  "Glitch", "glitch", false),
      (5,  "Candy Core Mine", "candy-core-mine", false),
      (6,  "Sweetie Plains", "sweetie-plains", false),
      (7,  "Peekaboo", "peekaboo", false),
      (8,  "Blitz Beach", "blitz-beach", false),
      (9,  "Long Nine", "long-nine", false),
      (10, "Sharkfood Island", "sharkfood-island", false),
      (11, "Acrospire", "acrospire", false),
      (12, "Thorn Gulley", "thorn-gulley", false),
      (13, "Reaper", "reaper", false),
      (14, "Skull D", "skull-duggery", false),
      (15, "War Garden", "war-garden", false),
      (16, "Sweet Tooth", "sweet-tooth", false),
      (17, "Sugar Rock", "sugar-rock", false),
      (18, "Mechanism", "mechanism", false);`,

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

		// The player identifier and name, their latest rankings in solos and duos.
		`CREATE TABLE "players" (
      "id"    INTEGER PRIMARY KEY,
      "guid"  TEXT NOT NULL UNIQUE,
      "name"  TEXT NOT NULL
    ) WITHOUT ROWID;`,

		`INSERT INTO players (id, guid, name) VALUES
      (0, "", "UNKNOWN");`,

		// The `solo_matches` metadata relates to an instance of a game between two
		// players.  This differs from the serialized replays that relate entirely
		// to each game's turns, or `solo_roles` which uniquely indexes the players
		// to their involvement in the match.
		`CREATE TABLE "solo_matches" (
      -- rowid INTEGER PRIMARY KEY is legacy "id"|Index field
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

		`CREATE TABLE "solo_roles" (
      -- rowid INTEGER PRIMARY KEY,
      "match_id" INTEGER NOT NULL,
      "player_id" INTEGER NOT NULL,
      "turn_order" INTEGER CHECK(turn_order > 0 AND turn_order <= 2),

      FOREIGN KEY (match_id) REFERENCES solo_matches (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_id) REFERENCES players (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,

      UNIQUE (match_id, turn_order) ON CONFLICT FAIL,
      UNIQUE (match_id, player_id) ON CONFLICT IGNORE
    );`,

		`CREATE TABLE "solo_standings" (
      -- rowid INTEGER PRIMARY KEY,
      "after" INTEGER NOT NULL UNIQUE,
      "until" INTEGER,

      "player_league" INTEGER NOT NULL,
      "player_rank"   INTEGER NOT NULL,
      "player_points" INTEGER DEFAULT 0,
      "player_delta"  INTEGER DEFAULT 0,

      "prev_points"   INTEGER DEFAULT 0,

      FOREIGN KEY (after) REFERENCES solo_roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (until) REFERENCES solo_roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_league) REFERENCES leagues (league_id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    );`,

		// The `duos_matches` metadata relates to an instance of a game between four
		// players.  This differs from the serialized replays that relate entirely
		// to each game's turns, or `duos_roles` which uniquely indexes the players
		// to their involvement in the match.
		`CREATE TABLE "duos_matches" (
      -- rowid INTEGER PRIMARY KEY is legacy "id"|Index field
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

		`CREATE TABLE "duos_roles" (
      -- rowid INTEGER PRIMARY KEY,
      "match_id" INTEGER NOT NULL,
      "player_id" INTEGER NOT NULL,
      "turn_order" INTEGER CHECK(turn_order > 0 AND turn_order <= 4),

      FOREIGN KEY (match_id) REFERENCES duos_matches (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (player_id) REFERENCES players (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,

      UNIQUE (match_id, turn_order) ON CONFLICT FAIL,
      UNIQUE (match_id, player_id) ON CONFLICT IGNORE
    );`,

		`CREATE TABLE "duos_standings" (
      -- rowid INTEGER PRIMARY KEY,
      "after" INTEGER NOT NULL UNIQUE,
      "until" INTEGER,

      "player_league" INTEGER NOT NULL,
      "player_rank"   INTEGER NOT NULL,
      "player_points" INTEGER DEFAULT 0,
      "player_delta"  INTEGER DEFAULT 0,

      "prev_points"   INTEGER DEFAULT 0,

      FOREIGN KEY (after) REFERENCES duos_roles (rowid)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (until) REFERENCES duos_roles (rowid)
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
