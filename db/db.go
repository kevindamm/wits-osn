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
// github:kevindamm/wits-osn/db/db.go

package db

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
type OsnDB interface {
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

// Opens a Sqlite db at indicated path and prepares queries.
// Does not create tables, only prepares the connection and statements.
//
// Asserts db exists and basic queries can be constructed.
// LOG(FATAL) on any error, database is returned unmodified.
func OpenOsnDB(filepath string) OsnDB {
	osndb, err := open_database(filepath)
	if err != nil {
		log.Fatal(err)
	}
	if err = osndb.PrepareQueries(); err != nil {
		log.Fatal(err)
	}
	return OsnDB(osndb)
}

// Opens a connection to the database but does not prepare any queries.
//
// Useful for getting a connection to create the required tables, Callers
// should prefer using [OpenOsnDB()], [CreateTables()], or methods on [OsnDB].
//
// Internally, this is used to get a connection without preparing queries (which
// would fail when required tables had not been created yet).
func open_database(filepath string) (*osndb, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	osndb := new(osndb)
	osndb.sqldb = db
	return osndb, nil
}

// Prepares [*sql.Stmt] statements for this database interface.
// (Automatically done if opening the database with [OpenOsnDB()] constructor.)
func (osndb *osndb) PrepareQueries() error {
	var err error = nil
	db := osndb.sqldb

	osndb.insertPlayer, err = db.Prepare(`INSERT INTO
	  players (id, name)
		VALUES (?, ?)`)
	if err != nil {
		return err
	}

	osndb.insertPlayerGCID, err = db.Prepare(`INSERT INTO
		player_gcid (player_id, gcid)
		VALUES (?, ?)`)
	if err != nil {
		return err
	}

	osndb.insertMatch, err = db.Prepare(`INSERT INTO
		matches (match_hash, competitive, season, start_time,
		  map_id, turn_count, version, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	osndb.insertPlayerRole, err = db.Prepare(`INSERT INTO
	  roles (match_id, player_id, turn_order)
		VALUES (?, ?, ?);`)
	if err != nil {
		return err
	}

	osndb.insertStanding, err = db.Prepare(`INSERT INTO
	  standings (after, player_league, player_rank, player_points, player_delta)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	osndb.updatePrevStanding, err = db.Prepare(`UPDATE standings
		SET until = ?
		WHERE after = ?;`)
	if err != nil {
		return err
	}

	osndb.selectMaps, err = db.Prepare(`SELECT
	  map_id, map_name, player_count, map_filename, map_theme, width, height
	  FROM maps
		WHERE NOT deprecated`)
	if err != nil {
		return err
	}

	osndb.selectPlayerByID, err = db.Prepare(`SELECT *
	  FROM players
	  WHERE id = ?;`)
	if err != nil {
		return err
	}

	osndb.selectPlayerByName, err = db.Prepare(`SELECT *
	  FROM players
	  WHERE name = ?;`)
	if err != nil {
		return err
	}

	osndb.selectMatches, err = db.Prepare(`SELECT *
		FROM matches
		WHERE match_hash = ?;`)
	if err != nil {
		return err
	}

	return nil
}

type osndb struct {
	sqldb *sql.DB

	cachedMaps    map[int]osn.Map
	cachedPlayers map[int]osn.Player

	insertPlayer       *sql.Stmt
	insertPlayerGCID   *sql.Stmt
	insertMatch        *sql.Stmt
	insertPlayerRole   *sql.Stmt
	insertStanding     *sql.Stmt
	updatePrevStanding *sql.Stmt

	selectMaps         *sql.Stmt
	selectPlayerByID   *sql.Stmt
	selectPlayerByName *sql.Stmt
	selectMatches      *sql.Stmt
}

func (db *osndb) Close() {
	db.sqldb.Close()
}

func (db *osndb) Player(id int) (osn.Player, error) {
	if db.cachedPlayers == nil {
		db.cachedPlayers = make(map[int]osn.Player)
	}
	if player, ok := db.cachedPlayers[id]; ok {
		return player, nil
	}

	row, err := db.selectPlayerByID.Query(id)
	if err != nil {
		return osn.UNKNOWN_PLAYER, err
	}
	if !row.Next() {
		return osn.UNKNOWN_PLAYER, fmt.Errorf("no player with ID %d", id)
	}

	player := osn.Player{}
	err = row.Scan(&player.ID.RowID, &player.Name)
	if err != nil {
		return player, err
	}

	// TODO also retrieve latest standings

	return player, nil
}

func (db *osndb) PlayerByName(name string) (osn.Player, error) {
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

func (db *osndb) InsertPlayer(player osn.Player) error {
	if _, ok := db.cachedPlayers[player.ID.RowID]; !ok {
		if db.cachedPlayers == nil {
			db.cachedPlayers = make(map[int]osn.Player)
		}
		_, err := db.insertPlayer.Exec(player.ID.RowID, player.Name)
		if err != nil {
			return err
		}
		db.cachedPlayers[player.ID.RowID] = player
	}

	if player.ID.GCID != "" {
		_, err := db.insertPlayerGCID.Exec(player.ID.RowID, player.ID.GCID)
		if err != nil {
			return err
		}
	}

	// TODO detect and insert standings

	return nil
}

func (db *osndb) AllMaps() ([]osn.Map, error) {
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

func (db *osndb) Map(id int) (osn.Map, error) {
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

func (db *osndb) MapByName(name string) (osn.Map, error) {
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

func (db *osndb) cachemaps() error {
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

func (db *osndb) Match(id string) (osn.LegacyMatch, error) {
	match := osn.UNKNOWN_MATCH

	return match, nil
}

func (db *osndb) InsertMatch(osn.LegacyMatch) error {
	// TODO
	return nil
}

func (db *osndb) InsertStanding(standing osn.PlayerStanding, previous int, current int) error {
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
