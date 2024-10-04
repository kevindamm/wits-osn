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
	"time"

	osn "github.com/kevindamm/wits-osn"
)

type LegacyMatchRecord struct {
	// Embedded so that type inference on Table[*LegacyMatchRecord] works.
	osn.LegacyMatch
}

func NewMatchRecord(match osn.LegacyMatch) *LegacyMatchRecord {
	return &LegacyMatchRecord{match}
}

func (*LegacyMatchRecord) Columns() []string {
	return []string{
		"rowid",
		"match_hash", "competitive", "season", "start_time",
		"map_id", "turn_count", "version", "fetch_status",
	}
}

func (record *LegacyMatchRecord) Values() ([]any, error) {
	return []any{
		record.MatchIndex,
		record.MatchHash,
		record.Competitive,
		record.Season,
		record.StartTime,
		record.MapID,
		record.TurnCount,
		record.Version,
		record.FetchStatus}, nil
}

func (record *LegacyMatchRecord) NamedValues() ([]driver.NamedValue, error) {
	return []driver.NamedValue{
		{
			Name:    "rowid",
			Ordinal: 0,
			Value:   record.MatchIndex},
		{
			Name:    "match_hash",
			Ordinal: 1,
			Value:   record.MatchHash},
		{
			Name:    "competitive",
			Ordinal: 2,
			Value:   record.Competitive},
		{
			Name:    "season",
			Ordinal: 3,
			Value:   record.Season},
		{
			Name:    "start_time",
			Ordinal: 4,
			Value:   record.StartTime},
		{
			Name:    "map_id",
			Ordinal: 5,
			Value:   record.MapID},
		{
			Name:    "turn_count",
			Ordinal: 6,
			Value:   record.TurnCount},
		{
			Name:    "version",
			Ordinal: 7,
			Value:   record.Version},
		{
			Name:    "fetch_status",
			Ordinal: 8,
			Value:   record.FetchStatus},
	}, nil
}

func (record *LegacyMatchRecord) ScanValues(values ...driver.Value) error {
	var ok bool
	record.MatchIndex, ok = values[0].(int64)
	if !ok {
		return fmt.Errorf("LegacyMatch.MatchIndex value %v not int64", values[0])
	}
	match_hash, ok := values[1].(string)
	if !ok {
		return fmt.Errorf("LegacyMatch.MatchHash value %v not string", values[1])
	}
	record.MatchHash = osn.GameID(match_hash)
	record.Competitive, ok = values[2].(bool)
	if !ok {
		return fmt.Errorf("LegacyMatch.Competitive value %v not boolean", values[2])
	}
	record.Season, ok = values[3].(int)
	if !ok {
		return fmt.Errorf("LegacyMatch.Season value %v not int", values[3])
	}
	record.StartTime, ok = values[4].(time.Time)
	if !ok {
		return fmt.Errorf("LegacyMatch.StartTime value %v not time.Time", values[4])
	}
	record.MapID, ok = values[5].(int)
	if !ok {
		return fmt.Errorf("LegacyMatch.MapID value %v not int", values[5])
	}
	record.TurnCount, ok = values[6].(int)
	if !ok {
		return fmt.Errorf("LegacyMatch.TurnCount value %v not int", values[6])
	}
	record.Version, ok = values[7].(int)
	if !ok {
		return fmt.Errorf("LegacyMatch.Version value %v not int", values[7])
	}
	record.FetchStatus, ok = values[8].(osn.FetchStatus)
	if !ok {
		return fmt.Errorf("LegacyMatch.FetchStatus value %v not valid", values[8])
	}
	return nil
}

func (record *LegacyMatchRecord) ScanRow(row *sql.Row) error {
	row.Scan(
		&record.MatchIndex,
		&record.MatchHash,
		&record.Competitive,
		&record.Season,
		&record.StartTime,
		&record.MapID,
		&record.TurnCount,
		&record.Version,
		&record.FetchStatus)
	return nil
}

func (record *LegacyMatchRecord) Scannables() []any {
	return []any{
		&record.MatchIndex,
		&record.MatchHash,
		&record.Competitive,
		&record.Season,
		&record.StartTime,
		&record.MapID,
		&record.TurnCount,
		&record.Version,
		&record.FetchStatus,
	}
}

// The `matches` metadata relates to an instance of a game between two
// players.  This differs from the serialized replays that define all of
// each game's turns, or `roles` which uniquely indexes the players to
// their involvement in the match.
type tableMatches struct {
	mutableBase[*LegacyMatchRecord]
	cached map[osn.GameID]osn.LegacyMatch
}

func MakeMatchesTable(sqldb *sql.DB) MutableTable[*LegacyMatchRecord] {
	return tableMatches{
		mutableBase[*LegacyMatchRecord]{tableBase[*LegacyMatchRecord]{
			sqldb:   sqldb,
			name:    "matches",
			NameCol: "match_hash"}},
		make(map[osn.GameID]osn.LegacyMatch)}
}

func (table tableMatches) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    -- rowid INTEGER PRIMARY KEY AUTOINCREMENT, -- legacy "id" or Index
    "match_hash"    TEXT NOT NULL UNIQUE,
    "competitive"   BOOLEAN,    -- league or friendly
    "season"        INTEGER,    -- seasons are of variable duration
    "start_time"    TIMESTAMP,  -- time at creation, UTC

    "map_id"        INTEGER,    -- MapEnum
    "turn_count"    INTEGER,    -- number of turns (= one ply) for the match

    "version"       INTEGER,    -- engine (runtime) version for this match
    "fetch_status"  INTEGER,    -- this match's fetch_status

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
type PlayerRoleRecord struct {
	MatchID   int64
	PlayerID  int64
	TurnOrder osn.PlayerColorEnum
}

func (*PlayerRoleRecord) Columns() []string {
	return []string{"match_id", "player_id", "turn_order"}
}

func (record *PlayerRoleRecord) Values() ([]any, error) {
	return []any{
			record.MatchID,
			record.PlayerID,
			record.TurnOrder},
		nil
}

func (record *PlayerRoleRecord) NamedValues() ([]driver.NamedValue, error) {
	return []driver.NamedValue{
		{
			Name:    "match_id",
			Ordinal: 1,
			Value:   record.MatchID},
		{
			Name:    "player_id",
			Ordinal: 2,
			Value:   record.PlayerID},
		{
			Name:    "turn_order",
			Ordinal: 3,
			Value:   record.TurnOrder},
	}, nil
}

func (record *PlayerRoleRecord) ScanValues(values ...driver.Value) error {
	var ok bool
	record.MatchID, ok = values[0].(int64)
	if !ok {
		return fmt.Errorf("PlayerRole.MatchID value %v not int64", values[0])
	}
	record.PlayerID, ok = values[1].(int64)
	if !ok {
		return fmt.Errorf("PlayerRole.PlayerID value %v not int64", values[1])
	}
	record.TurnOrder, ok = values[2].(osn.PlayerColorEnum)
	if !ok {
		return fmt.Errorf("PlayerRole.TurnOrder value %v not 1-4", values[2])
	}
	return nil
}

func (record *PlayerRoleRecord) ScanRow(row *sql.Row) error {
	return row.Scan(&record.MatchID, &record.PlayerID, &record.TurnOrder)
}

func (record *PlayerRoleRecord) Scannables() []any {
	return []any{
		&record.MatchID,
		&record.PlayerID,
		&record.TurnOrder}
}

type tableRoles struct {
	mutableBase[*PlayerRoleRecord]
}

func MakeRolesTable(sqldb *sql.DB) MutableTable[*PlayerRoleRecord] {
	return tableRoles{
		mutableBase[*PlayerRoleRecord]{tableBase[*PlayerRoleRecord]{
			sqldb: sqldb,
			name:  "roles"}}}
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
