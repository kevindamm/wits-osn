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
// github:kevindamm/wits-osn/cmd/fetch/db_test.go

package main_test

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	osn "github.com/kevindamm/wits-osn"
	main "github.com/kevindamm/wits-osn/cmd/fetch"
)

func TestDB(t *testing.T) {
	file, err := os.CreateTemp("", "wits-*.db")
	if err != nil {
		t.Errorf("error creating temporary file for DB: %s\n", err)
	}
	file.Close()

	create_tables := true
	var db main.WitsDB = main.OpenWitsDB(file.Name(), create_tables)
	defer db.Close()

	maps, err := db.Maps()
	if err != nil {
		t.Errorf("MapNames() error: %s\n", err)
	}

	if len(maps) != 17 {
		t.Errorf("retrieved unexpected number of maps %d", len(maps))
		t.Log(maps)
	}

	err = db.InsertPlayer(osn.Player{
		ID: osn.PlayerID{RowID: 1, GCID: "abcde"}, Name: "First"})
	if err != nil {
		t.Error(err)
	}
	err = db.InsertPlayer(osn.Player{
		ID: osn.PlayerID{RowID: 2, GCID: "bcdef"}, Name: "2nd"})
	if err != nil {
		t.Error(err)
	}

	check_player(t, db, 1, "First")
	check_player(t, db, 2, "2nd")
}

func check_player(t *testing.T, db main.WitsDB, id int64, name string) {
	player, err := db.Player(id)
	if err != nil {
		t.Errorf("error retrieving player (id=1)\n%s", err)
	}
	if player.Name != name {
		t.Errorf("incorrect name for player (id=1): %s", player.Name)
	}

	player, err = db.PlayerByName(name)
	if err != nil {
		t.Errorf("error retrieving player (name=First)\n%s", err)
	}
	if player.ID.RowID != id {
		t.Errorf("incorrect ID for player (name=First): %d", player.ID.RowID)
	}
}
