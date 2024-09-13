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
// github:kevindamm/wits-osn/cmd/fetch/db.go

package main

import (
	"database/sql"
	"fmt"
	"log"

	osn "github.com/kevindamm/wits-osn"
	_ "github.com/mattn/go-sqlite3"
)

// Interface for simplifying the interaction with a backing database.
//
// DB includes match metadata, map identities, player history and replay index.
type WitsDB interface {
	Close() // also closes the SqlDB

	InsertPlayer(osn.Player) error
	InsertMatch(osn.LegacyMatch) error
	InsertStanding(osn.PlayerStanding, int, int) error

	AllMaps() ([]osn.Map, error)
	Map(id int) (osn.Map, error)
	MapByName(name string) (osn.Map, error)

	Player(id int) (osn.Player, error)
	PlayerByName(name string) (osn.Player, error)

	Match(id string) (osn.LegacyMatch, error)
}

func assertNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Opens a Sqlite db at indicated path and prepares queries.
// Does not create tables, only prepares the connection and statements.
//
// Asserts db exists and basic queries can be constructed.
// LOG(FATAL) on any error but database is not modified.
func OpenWitsDB(filepath string) WitsDB {
	witsdb, err := open_database(filepath)
	if err != nil {
		log.Fatal(err)
	}
	if err = witsdb.PrepareQueries(); err != nil {
		log.Fatal(err)
	}
	return WitsDB(witsdb)
}

// Opens a connection to the database but does not prepare any queries.
//
// Useful for getting a connection to create the required tables,  Callers from
// other packages should use OpenWitsDB() or CreateTables(), & WitsDB interface.
func open_database(filepath string) (*witsdb, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	witsdb := new(witsdb)
	witsdb.sqldb = db
	return witsdb, nil
}

// Prepares [*sql.Stmt] statements for this database interface.
// (Automatically done if opening the database with [OpenWitsDB()] constructor.)
func (witsdb *witsdb) PrepareQueries() error {
	var err error = nil
	db := witsdb.sqldb

	witsdb.insertPlayer, err = db.Prepare(`INSERT INTO
	  players (id, name)
		VALUES (?, ?)`)
	if err != nil {
		return err
	}

	witsdb.insertMatch, err = db.Prepare(`INSERT INTO
		matches (rowid, match_hash, competitive, season, start_time,
		  map_id, turn_count, version, osn_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	witsdb.insertPlayerRole, err = db.Prepare(`INSERT INTO
	  roles (match_id, player_id, turn_order)
		VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}

	witsdb.insertStanding, err = db.Prepare(`INSERT INTO
	  standings (after, player_league, player_rank, player_points, player_delta)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	witsdb.updatePrevStanding, err = db.Prepare(`UPDATE standings
		SET until = ?
		WHERE after = ?;`)
	if err != nil {
		return err
	}

	witsdb.selectMaps, err = db.Prepare(`SELECT
	  map_id, map_name, player_count, map_filename, map_theme, width, height
	  FROM maps
		WHERE NOT deprecated`)
	if err != nil {
		return err
	}

	witsdb.selectPlayerByID, err = db.Prepare(`SELECT *
	  FROM players
	  WHERE id = ?;`)
	if err != nil {
		return err
	}

	witsdb.selectPlayerByName, err = db.Prepare(`SELECT *
	  FROM players
	  WHERE name = ?;`)
	if err != nil {
		return err
	}

	witsdb.selectMatches, err = db.Prepare(`SELECT *
		FROM matches
		WHERE match_hash = ?;`)
	if err != nil {
		return err
	}

	return nil
}

type witsdb struct {
	sqldb *sql.DB

	cachedMaps    map[int]osn.Map
	cachedPlayers map[int]osn.Player

	insertPlayer       *sql.Stmt
	insertMatch        *sql.Stmt
	insertPlayerRole   *sql.Stmt
	insertStanding     *sql.Stmt
	updatePrevStanding *sql.Stmt

	selectMaps         *sql.Stmt
	selectPlayerByID   *sql.Stmt
	selectPlayerByName *sql.Stmt
	selectMatches      *sql.Stmt
}

func (db *witsdb) Close() {
	db.sqldb.Close()
}

func (db *witsdb) Player(id int) (osn.Player, error) {
	player := osn.UNKNOWN_PLAYER
	row, err := db.selectPlayerByID.Query(id)
	if err != nil {
		return player, err
	}
	if !row.Next() {
		return player, fmt.Errorf("no player with ID %d", id)
	}

	err = row.Scan(&player.ID.RowID, &player.Name)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (db *witsdb) PlayerByName(name string) (osn.Player, error) {
	player := osn.UNKNOWN_PLAYER
	row, err := db.selectPlayerByName.Query(name)
	if err != nil {
		return player, err
	}
	if !row.Next() {
		return player, fmt.Errorf("no player with name %s", name)
	}

	err = row.Scan(&player.ID.RowID, &player.Name)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (db *witsdb) InsertPlayer(player osn.Player) error {
	if _, ok := db.cachedPlayers[player.ID.RowID]; ok {
		// don't need to insert the same player twice
		return nil
	}

	_, err := db.insertPlayer.Exec(player.ID.RowID, player.Name)
	if err != nil {
		db.cachedPlayers[player.ID.RowID] = player
	}
	return err
}

func (db *witsdb) AllMaps() ([]osn.Map, error) {
	maps := make([]osn.Map, 0)
	if rows, err := db.selectMaps.Query(); err != nil {
		return maps, err
	} else {
		defer rows.Close()

		mapobj := osn.Map{}
		for rows.Next() {
			// map_id, map_name, player_count, map_filename, map_theme, width, height
			rows.Scan(
				&mapobj.MapID, &mapobj.Name, &mapobj.PlayerCount,
				&mapobj.Filename, &mapobj.Theme,
				&mapobj.Width, &mapobj.Height)
			maps = append(maps, mapobj)
		}
	}

	return maps, nil
}

func (db *witsdb) Map(id int) (osn.Map, error) {
	if len(db.cachedMaps) == 0 {
		if err := db.cachemaps(); err != nil {
			return osn.UNKNOWN_MAP, err
		}
	}

	mapobj, found := db.cachedMaps[id]
	if !found {
		return osn.UNKNOWN_MAP, fmt.Errorf("unrecognized map ID %d", id)
	}
	return mapobj, nil
}

func (db *witsdb) MapByName(name string) (osn.Map, error) {
	if len(db.cachedMaps) == 0 {
		if err := db.cachemaps(); err != nil {
			return osn.UNKNOWN_MAP, err
		}
	}
	for _, mapobj := range db.cachedMaps {
		if mapobj.Name == name {
			return mapobj, nil
		}
	}

	return osn.UNKNOWN_MAP, fmt.Errorf("unknown map named %s", name)
}

func (db *witsdb) cachemaps() error {
	maps, err := db.AllMaps()
	db.cachedMaps = make(map[int]osn.Map)
	if err != nil {
		return err
	}
	for _, mapobj := range maps {
		db.cachedMaps[int(mapobj.MapID)] = mapobj
	}
	return nil
}

func (db *witsdb) Match(id string) (osn.LegacyMatch, error) {
	match := osn.UNKNOWN_MATCH

	return match, nil
}

func exec(db *sql.DB, query string) error {
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

func (db *witsdb) InsertMatch(osn.LegacyMatch) error {
	// TODO

	return nil
}

func (db *witsdb) InsertStanding(standing osn.PlayerStanding, previous int, current int) error {
	if _, err := db.insertStanding.Exec(
		current, standing.League(), standing.Rank(),
		standing.PointsAfter(), standing.Delta()); err != nil {
		return err
	}
	if _, err := db.updatePrevStanding.Exec(current, previous); err != nil {
		return err
	}

	return nil
}
