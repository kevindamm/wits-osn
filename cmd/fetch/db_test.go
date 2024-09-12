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
	file, err := os.CreateTemp("", "osnwits-*.db")
	if err != nil {
		t.Errorf("error creating temporary file for DB: %s\n", err)
	}
	file.Close()

	// will also create the database that doesn't exist yet
	err = main.CreateTablesAndClose(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Closing and re-opening the database should still include the tables above.
	db := main.OpenWitsDB(file.Name())

	maps, err := db.AllMaps()
	if err != nil {
		t.Errorf("MapNames() error: %s\n", err)
	}
	t.Logf("maps: %v", maps)
	if len(maps) != 17 {
		t.Errorf("retrieved unexpected number of maps %d", len(maps))
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

	metadata := osn.LegacyReplayMetadata{
		Index:       "5",
		GameID:      "ag5vdXR3aXR0ZXJzZ2FtZXIQCxIIR2FtZVJvb20Y9-5HDA",
		GameType:    "2",
		LeagueMatch: "1",
		Created:     "",
		Season:      "1",

		WitsVersion: "1000",
		MapID:       "7",
		MapName:     "Peekaboo",
		MapTheme:    "2",

		TurnCount: "25",
		ViewCount: "35",
		LikeCount: "1",

		Player1_ID:     "2",
		Player1_Name:   "Alvendor",
		Player1_League: "5",
		Player1_Race:   "3",
		Player1_Wins:   "0",
		Player1_BaseHP: "0",

		Player2_ID:     "3",
		Player2_Name:   "Lenoxe",
		Player2_League: "5",
		Player2_Race:   "3",
		Player2_Wins:   "1",
		Player2_BaseHP: "5",

		FirstPlayer: "3",
	}
	match := metadata.ToLegacyMatch()

	err = db.InsertMatch(match)
	if err != nil {
		t.Errorf("error when inserting match metadata:\n%s", err)
	}

	check_match(t, db, match)
}

func check_player(t *testing.T, db main.WitsDB, id int, name string) {
	player, err := db.Player(id)
	if err != nil {
		t.Errorf("error retrieving player (id=%d)\n%s", id, err)
	}
	if player.Name != name {
		t.Errorf("incorrect name for player (id=%d): %s", id, player.Name)
	}

	player, err = db.PlayerByName(name)
	if err != nil {
		t.Errorf("error retrieving player (name=%s)\n%s", name, err)
	}
	if player.ID.RowID != id {
		t.Errorf("incorrect ID for player (name=%s): %d", name, player.ID.RowID)
	}
}

func check_match(t *testing.T, db main.WitsDB, expected osn.LegacyMatch) {
	match, err := db.Match(expected.OsnMatchID)
	if err != nil {
		t.Errorf("error retrieving match %s\n", match.OsnMatchID)
	}

	if match.OsnMatchID != expected.OsnMatchID {

	}
}
