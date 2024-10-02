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

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	osn "github.com/kevindamm/wits-osn"
)

type LegacyMatchRecord osn.LegacyMatch

func (*LegacyMatchRecord) Columns() []string {
	return []string{
		"match_hash", "competitive", "season", "start_time",
		"map_id", "turn_count", "version", "status",
	}
}

func (record *LegacyMatchRecord) Values() ([]driver.Value, error) {

	// TODO
	return []driver.Value{}, nil
}

func (record *LegacyMatchRecord) NamedValues() ([]driver.NamedValue, error) {

	// TODO
	return []driver.NamedValue{}, nil
}

func (record *LegacyMatchRecord) ScanValues(...driver.Value) error {

	// TODO
	return nil
}

func (record *LegacyMatchRecord) ScanRow(*sql.Row) error {

	// TODO
	return nil
}

// The `matches` metadata relates to an instance of a game between two
// players.  This differs from the serialized replays that define all of
// each game's turns, or `roles` which uniquely indexes the players to
// their involvement in the match.
type tableMatches struct {
	tableBase[*LegacyMatchRecord]
	cached map[osn.GameID]osn.LegacyMatch
}

func MakeMatchesTable(sqldb *sql.DB) Table[*LegacyMatchRecord] {
	return tableMatches{
		tableBase[*LegacyMatchRecord]{
			sqldb:   sqldb,
			name:    "matches",
			NameCol: "match_hash"},
		make(map[osn.GameID]osn.LegacyMatch)}
}

func (table tableMatches) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
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
  );`, table.name)
}

func (tableMatches) SqlInit() string {
	return `CREATE UNIQUE INDEX match_hashes ON matches (match_hash);`
}

// Relation for which players are participating in which matches, and the turn
// order they are assigned to.  Appropriate for both 1v1 and 2v2 matches.
type PlayerRoleRecord osn.PlayerRole

func (PlayerRoleRecord) Columns() []string {
	return []string{"match_id", "player_id", "turn_order"}
}

func (record PlayerRoleRecord) Values() ([]driver.Value, error) {

	// TODO
	return []driver.Value{}, nil
}

func (record PlayerRoleRecord) NamedValues() ([]driver.NamedValue, error) {

	// TODO
	return nil, nil
}

func (record PlayerRoleRecord) ScanValues(...driver.Value) error {

	// TODO
	return nil
}

func (record PlayerRoleRecord) ScanRow(*sql.Row) error {

	// TODO
	return nil
}

type tableRoles struct {
	tableBase[*PlayerRoleRecord]
}

func MakeRolesTable(sqldb *sql.DB) Table[*PlayerRoleRecord] {
	return tableRoles{
		tableBase[*PlayerRoleRecord]{
			sqldb: sqldb,
			name:  "roles"}}
}

func (table tableRoles) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
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
    );`, table.name)
}

func (table tableRoles) SqlInit() string {
	return ""
}
