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
// github:kevindamm/wits-osn/db/players.go

package db

import (
	"database/sql"
	"database/sql/driver"

	osn "github.com/kevindamm/wits-osn"
)

// This satisfies the Record interface for use as Table[T]
// while being only a thin wrapper around the Player struct.
type PlayerRecord osn.Player

var PlayerSchema = []string{
	// The player identifier and in-game display name,
	// used for foreign key relations in replay-related tables.
	`CREATE TABLE "players" (
      "id"    INTEGER PRIMARY KEY,
      "gcid"  TEXT UNIQUE,
      "name"  TEXT NOT NULL
    ) WITHOUT ROWID;`,

	`INSERT INTO players (id, gcid, name) VALUES (0, NULL, "UNKNOWN");`,

	// An index is created to assist lookup of players by name.
	`CREATE UNIQUE INDEX player_names
      ON players (name);`,
}

func (player PlayerRecord) ScanRecord(row *sql.Row) error {

	return nil
}

func (player *PlayerRecord) RecordValues() ([]driver.Value, error) {

	return []driver.Value{}, nil
}

var PlayerRoleSchema = []string{
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
}

var PlayerStandingsSchema = []string{
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
}
