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

package db_test

import (
	"log"
	"os"
	"testing"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	file, err := os.CreateTemp("", "osn-*.db")
	if err != nil {
		t.Errorf("error creating temporary file for DB: %s\n", err)
	}
	file.Close()

	// will also create the database that doesn't exist yet
	err = db.CreateTablesAndClose(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Closing and re-opening the database should still include the tables above.
	db := db.OpenOsnDB(file.Name())

	mapobj, err := db.Map(1)
	if err != nil {
		t.Errorf("could not find map ID 1: %s", err)
	}
	if mapobj.MapID != 1 || mapobj.Name != "Machination" {
		t.Error("retrieved incorrect map for ID 1")
	}

	mapobj, err = db.Map(2)
	if err == nil || mapobj.MapID != 0 {
		t.Error("deprecated map should not be included")
	}

	maps, err := db.AllMaps()
	if err != nil {
		t.Errorf("MapNames() error: %s\n", err)
	} else {
		t.Logf("maps: %v", maps)
		if len(maps) != 17 {
			t.Errorf("retrieved unexpected number of maps %d", len(maps))
		}
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
		GameID:      "ag5vdXR3aXR0ZXJzZ2FtZXIQCxIIR2FtZVJvb20Y9-5HDA",
		NumPlayers:  "2",
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
	log.Print(metadata)

	match := metadata.ToLegacyMatch()

	err = db.InsertMatch(match)
	if err != nil {
		t.Errorf("error when inserting match metadata:\n%s", err)
	}

	check_match(t, db, match)
}

func check_player(t *testing.T, db db.OsnDB, id int, name string) {
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

func check_match(t *testing.T, db db.OsnDB, expected osn.LegacyMatch) {
	match, err := db.Match(string(expected.MatchID))
	if err != nil {
		t.Errorf("error retrieving match %s\n", match.MatchID)
	}

	if match.MatchID != expected.MatchID {
		t.Errorf("match ID incorrect; got %s (expected %s)",
			match.MatchID, expected.MatchID)
	}
	if match.Competitive != expected.Competitive {
		t.Errorf("competitive %v != expected.competitive %v", match.Competitive, expected.Competitive)
	}
	if match.Season != expected.Season {
		t.Errorf("season %d != expected.season %d", match.Season, expected.Season)
	}
	if match.StartTime != expected.StartTime {
		t.Errorf("start_time %s != expected.start_time %s",
			match.StartTime, expected.StartTime)
	}
	if match.MapID != expected.MapID {
		t.Errorf("map_id %d != expected.map_id %d", match.MapID, expected.MapID)
	}
	if match.TurnCount != expected.TurnCount {
		t.Errorf("turn_count %d != expected.turn_count %d", match.TurnCount, expected.TurnCount)
	}

}