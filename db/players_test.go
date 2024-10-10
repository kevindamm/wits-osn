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
// github:kevindamm/wits-osn/db/players_test.go

package db_test

import (
	"testing"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

func TestPlayersTable(t *testing.T) {
	osndb := db.OpenOsnDB(":memory:")
	osndb.MustCreateAndPopulateTables()

	record := db.NewPlayerRecord()
	record.RowID = 1
	record.GCID = ""
	record.Name = "Player1"
	osndb.Players().Insert(record)

	player, err := osndb.Players().Get(1)
	if err != nil {
		t.Error(err)
	}
	if player.Name != "Player1" {
		t.Errorf("player name %s expected Player1", player.Name)
	}

	err = osndb.Players().Insert(&db.PlayerRecord{osn.Player{
		RowID: 1, GCID: "abcde", Name: "First"}})
	if err != nil {
		t.Error(err)
	}
	err = osndb.Players().Insert(&db.PlayerRecord{osn.Player{
		RowID: 2, GCID: "bcdef", Name: "2nd"}})
	if err != nil {
		t.Error(err)
	}

	check_player(t, osndb, 1, "First")
	check_player(t, osndb, 2, "2nd")

}

func check_player(t *testing.T, db db.OsnDB, id int64, name string) {
	player, err := db.Players().Get(id)
	if err != nil {
		t.Errorf("error retrieving player (id=%d)\n%s", id, err)
	}
	if player.Name != name {
		t.Errorf("incorrect name for player (id=%d): %s", id, player.Name)
	}

	player, err = db.Players().GetByName(name)
	if err != nil {
		t.Errorf("error retrieving player (name=%s)\n%s", name, err)
	}
	if player.RowID != id {
		t.Errorf("incorrect ID for player (name=%s): %d", name, player.RowID)
	}
}
