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
	Close()

	InsertPlayer(osn.Player) error
	InsertMatch(osn.LegacyReplayMetadata) error
	//InsertStanding(osn.PlayerStanding) error

	Maps() ([]osn.Map, error)
	Player(id int64) (osn.Player, error)
	PlayerByName(name string) (osn.Player, error)
}

func assertNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenWitsDB(filepath string, create_tables bool) WitsDB {
	db, err := sql.Open("sqlite3", filepath)
	assertNil(err)

	witsdb := new(witsdb)
	witsdb.sqldb = db

	if create_tables {
		CreateTables(witsdb)
	}

	witsdb.insertPlayer, err = db.Prepare(`INSERT INTO
	  players (id, guid, name) VALUES (?, ?, ?)`)
	assertNil(err)

	witsdb.selectEnums, err = db.Prepare(`SELECT * FROM maps WHERE NOT deprecated`)
	assertNil(err)

	witsdb.selectPlayerByID, err = db.Prepare(`SELECT * FROM players
	  WHERE id = ?;`)
	assertNil(err)

	witsdb.selectPlayerByName, err = db.Prepare(`SELECT * FROM players
	  WHERE name = ?;`)
	assertNil(err)

	return witsdb
}

type witsdb struct {
	sqldb *sql.DB

	insertPlayer     *sql.Stmt
	insertMatch      *sql.Stmt
	insertPlayerRole *sql.Stmt

	selectEnums        *sql.Stmt
	selectPlayerByID   *sql.Stmt
	selectPlayerByName *sql.Stmt
	selectMatches      *sql.Stmt
}

func (db *witsdb) Close() {
	db.sqldb.Close()
}

func (db *witsdb) InsertPlayer(player osn.Player) error {
	_, err := db.insertPlayer.Exec(
		player.ID.RowID, player.ID.GCID, player.Name)
	return err
}

func (db *witsdb) InsertMatch(osn.LegacyReplayMetadata) error {
	// TODO

	return nil
}

//func (db *witsdb) AddReplay(osn.PlayerStanding) error {
//	return nil
//}

func (db *witsdb) Maps() ([]osn.Map, error) {
	maps := make([]osn.Map, 0)
	rows, err := db.selectEnums.Query()
	if err != nil {
		return maps, err
	}
	defer rows.Close()

	mapobj := osn.Map{}
	for rows.Next() {
		rows.Scan(&mapobj.MapID, &mapobj.Name)
		maps = append(maps, mapobj)
	}

	return maps, nil
}

func (db *witsdb) Player(id int64) (osn.Player, error) {
	player := osn.UnknownPlayer()
	row, err := db.selectPlayerByID.Query(id)
	if err != nil {
		return player, err
	}
	if !row.Next() {
		return player, fmt.Errorf("no player with ID %d", id)
	}

	err = row.Scan(&player.ID.RowID, &player.ID.GCID, &player.Name)
	if err != nil {
		return player, err
	}
	return player, nil
}

func (db *witsdb) PlayerByName(name string) (osn.Player, error) {
	player := osn.UnknownPlayer()
	row, err := db.selectPlayerByName.Query(name)
	if err != nil {
		return player, err
	}
	if !row.Next() {
		return player, fmt.Errorf("no player with name %s", name)
	}

	err = row.Scan(&player.ID.RowID, &player.ID.GCID, &player.Name)
	if err != nil {
		return player, err
	}
	return player, nil
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
