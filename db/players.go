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

// This satisfies the Record interface for use as Table[T]
// while being only a thin wrapper around the Player struct.
type PlayerRecord osn.Player

func (player PlayerRecord) Columns() []string {
	return []string{
		"id",
		"gcid",
		"name",
	}
}

func (player *PlayerRecord) Values() ([]driver.Value, error) {

	// TODO
	return []driver.Value{}, nil
}

func (player *PlayerRecord) NamedValues() ([]driver.NamedValue, error) {

	// TODO
	return nil, nil
}

func (player *PlayerRecord) ScanValues(values ...driver.Value) error {

	// TODO
	return nil
}

func (player *PlayerRecord) ScanRow(row *sql.Row) error {
	// TODO
	return nil
}

type tablePlayers struct {
	tableBase[*PlayerRecord]
	cachedPlayers map[int64]*PlayerRecord
}

func MakePlayersTable(sqldb *sql.DB) Table[*PlayerRecord] {
	return tablePlayers{
		tableBase: tableBase[*PlayerRecord]{
			sqldb:   sqldb,
			name:    "players",
			Primary: "id",
			NameCol: "name",
		},
		cachedPlayers: make(map[int64]*PlayerRecord)}
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

func (record *StandingsRecord) Values() ([]driver.Value, error) {
	// TODO
	return []driver.Value{
		record.After,
		record.Until,
		record.League(),
		record.Rank(),
		record.PointsAfter(),
		record.Delta(),
	}, nil
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
			Value:   record.League(),
		},
		{
			Name:    "player_rank",
			Ordinal: 4,
			Value:   record.Rank(),
		},
		{
			Name:    "player_points",
			Ordinal: 5,
			Value:   record.PointsAfter(),
		},
		{
			Name:    "player_delta",
			Ordinal: 6,
			Value:   record.Delta(),
		},
	}, nil
}

func (record *StandingsRecord) ScanValues(...driver.Value) error {
	// TODO
	return nil
}

func (record *StandingsRecord) ScanRow(row *sql.Row) error {
	// TODO
	return nil
}

type tableStandings struct {
	tableBase[*StandingsRecord]
	cached map[int64]*StandingsRecord
}

func MakeStandingsTable(sqldb *sql.DB) Table[*StandingsRecord] {
	return tableStandings{
		tableBase[*StandingsRecord]{
			sqldb: sqldb,
			name:  "standings"},
		make(map[int64]*StandingsRecord)}
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
