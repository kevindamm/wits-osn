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
// github:kevindamm/wits-osn/db/matches.go

package db

var MatchesSchema = []string{
	// The `matches` metadata relates to an instance of a game between two
	// players.  This differs from the serialized replays that relate entirely
	// to each game's turns, or `solo_roles` which uniquely indexes the players
	// to their involvement in the match.
	`CREATE TABLE "matches" (
      -- rowid INTEGER PRIMARY KEY AUTOINCREMENT, -- legacy "id" or Index
      "match_hash"   TEXT NOT NULL UNIQUE,
      "competitive"  BOOLEAN,    -- league or friendly
      "season"       INTEGER,    -- seasons are of variable duration
      "start_time"   TIMESTAMP,  -- time at creation, UTC

      "map_id"       INTEGER,    -- MapEnum
      "turn_count"   INTEGER,    -- number of turns (= one ply) for the match

      "version"      INTEGER,    -- engine (runtime) version for this match
      "status"       INTEGER,    -- this match's fetch_status

      FOREIGN KEY (map_id)
        REFERENCES maps (map_id)
        ON DELETE CASCADE ON UPDATE NO ACTION,
      FOREIGN KEY (status)
        REFERENCES fetch_status (id)
        ON DELETE CASCADE ON UPDATE NO ACTION
    );`,

	`CREATE UNIQUE INDEX match_hashes
      ON matches (match_hash);`,
}
