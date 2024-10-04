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
	"fmt"

	osn "github.com/kevindamm/wits-osn"
)

// This satisfies the Record interface for use as Table[*Record]
// while being only a thin wrapper around the Player struct.
type PlayerRecord struct {
	// For it to be a pointer-receiver in the following methods, it needs to be an
	// embedded struct member rather than a simple alias type.
	osn.Player
}

func (player *PlayerRecord) Columns() []string {
	return []string{
		"id",
		"gcid",
		"name",
	}
}

func (player *PlayerRecord) Values() ([]any, error) {
	return []any{
			player.RowID,
			player.GCID,
			player.Name},
		nil
}

func (player *PlayerRecord) NamedValues() ([]driver.NamedValue, error) {
	return []driver.NamedValue{
		{
			Name:    "id",
			Ordinal: 0,
			Value:   player.RowID},
		{
			Name:    "gcid",
			Ordinal: 1,
			Value:   player.GCID},
		{
			Name:    "name",
			Ordinal: 2,
			Value:   player.Name},
	}, nil
}

func (player *PlayerRecord) ScanValues(values ...driver.Value) error {
	var ok bool
	player.RowID, ok = values[0].(int64)
	if !ok {
		return fmt.Errorf("Player.RowID value %v not int64", values[0])
	}
	player.GCID, ok = values[1].(string)
	if !ok {
		return fmt.Errorf("Player.GCID value %v not string", values[1])
	}
	player.Name, ok = values[2].(string)
	if !ok {
		return fmt.Errorf("Player.Name value %v not string", values[2])
	}
	return nil
}

func (player *PlayerRecord) ScanRow(row *sql.Row) error {
	return row.Scan(&player.RowID, &player.GCID, &player.Name)
}

func (player *PlayerRecord) Scannables() []any {
	return []any{
		&player.RowID,
		&player.GCID,
		&player.Name}
}

type tablePlayers struct {
	mutableBase[*PlayerRecord]
	cachedPlayers map[int64]*osn.Player
}

func MakePlayersTable(sqldb *sql.DB) MutableTable[*PlayerRecord] {
	return tablePlayers{
		mutableBase: mutableBase[*PlayerRecord]{tableBase[*PlayerRecord]{
			sqldb:   sqldb,
			name:    "players",
			Primary: "id",
			NameCol: "name"}},
		cachedPlayers: make(map[int64]*osn.Player)}
}

func (table tablePlayers) SqlCreate() string {
	return fmt.Sprintf(`CREATE TABLE "%s" (
    "id"    INTEGER PRIMARY KEY,
    "gcid"  TEXT UNIQUE,
    "name"  TEXT NOT NULL
  ) WITHOUT ROWID;`, table.name)
}

func (table tablePlayers) SqlInit() string {
	return fmt.Sprintf(`
  INSERT INTO %s (id, gcid, name) VALUES (0, NULL, "UNKNOWN");
	
  CREATE UNIQUE INDEX player_names ON %s (name);`,
		table.name, table.name)
}

type StandingsRecord struct {
	After int64
	Until int64
	osn.PlayerStanding
}

func (*StandingsRecord) Columns() []string {
	return []string{
		"after_role",
		"until_role",
		"player_league",
		"player_rank",
		"player_points",
		"player_delta",
	}
}

func (record *StandingsRecord) Values() ([]any, error) {
	return []any{
			record.After,
			record.Until,
			record.League,
			record.Rank,
			record.Points,
			record.Delta},
		nil
}

func (record *StandingsRecord) NamedValues() ([]driver.NamedValue, error) {
	return []driver.NamedValue{
		{
			Name:    "after_role",
			Ordinal: 1,
			Value:   record.After,
		},
		{
			Name:    "until_role",
			Ordinal: 2,
			Value:   record.Until,
		},
		{
			Name:    "player_league",
			Ordinal: 3,
			Value:   record.League,
		},
		{
			Name:    "player_rank",
			Ordinal: 4,
			Value:   record.Rank,
		},
		{
			Name:    "player_points",
			Ordinal: 5,
			Value:   record.Points,
		},
		{
			Name:    "player_delta",
			Ordinal: 6,
			Value:   record.Delta,
		},
	}, nil
}

func (record *StandingsRecord) ScanValues(values ...driver.Value) error {
	var ok bool
	record.After, ok = values[0].(int64)
	if !ok {
		return fmt.Errorf("Standings.After value %v not int64", values[0])
	}
	record.Until, ok = values[1].(int64)
	if !ok {
		return fmt.Errorf("Standings.Before value %v not int64", values[1])
	}
	league, ok := values[2].(uint8)
	if !ok {
		return fmt.Errorf("Standings.League value %v not uint8", values[2])
	}
	record.League = osn.LeagueEnum(league)
	rank, ok := values[3].(uint8)
	if !ok {
		return fmt.Errorf("Standings.Rank value %v not uint8", values[3])
	}
	record.Rank = osn.LeagueRank(rank)
	record.Points, ok = values[4].(uint16)
	if !ok {
		return fmt.Errorf("Standings.Points value %v not uint16", values[4])
	}
	record.Delta, ok = values[5].(int8)
	if !ok {
		return fmt.Errorf("Standings.Delta value %v not int8", values[5])
	}
	return nil
}

func (record *StandingsRecord) ScanRow(row *sql.Row) error {
	return row.Scan(record.Scannables())
}

func (record *StandingsRecord) Scannables() []any {
	return []any{
		&record.After,
		&record.Until,
		&record.League,
		&record.Rank,
		&record.Points,
		&record.Delta}
}

type tableStandings struct {
	mutableBase[*StandingsRecord]
	cached map[int64]osn.PlayerStanding
}

func MakeStandingsTable(sqldb *sql.DB) MutableTable[*StandingsRecord] {
	return tableStandings{
		mutableBase[*StandingsRecord]{tableBase[*StandingsRecord]{
			sqldb: sqldb,
			name:  "standings"}},
		make(map[int64]osn.PlayerStanding)}
}

func (tableStandings) SqlCreate() string {
	return `CREATE TABLE "standings" (
    -- rowid INTEGER PRIMARY KEY,
    "after_role" INTEGER NOT NULL UNIQUE,
    "until_role" INTEGER,

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
  );`
}

func (tableStandings) SqlInit() string {
	return ""
}
