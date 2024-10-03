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
// github:kevindamm/wits-osn/db/db_test.go

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
	osndb := db.OpenOsnDB(file.Name())

	mapobj, err := osndb.Map("machination")
	if err != nil {
		t.Errorf("could not find map ID 1: %s", err)
	}
	if mapobj.MapID != 1 || mapobj.Name != "Machination" {
		t.Error("retrieved incorrect map for ID 1")
	}

	mapobj, err = osndb.Map("foundry")
	if err != nil {
		t.Errorf("could not find map ID 3 (foundry): %s", err)
	}
	if mapobj.MapID != 3 || mapobj.Name != "Foundry" {
		t.Error("retrieved incorrect map for ID 3")
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

	metadata := osn.LegacyReplayMetadata{
		GameID:      "ag5vdXR3aXR0ZXJzZ2FtZXIQCxIIR2FtZVJvb20Y9-5HDA",
		NumPlayers:  "2",
		LeagueMatch: "1",
		Created:     "2012-08-05 15:14:31",
		Season:      "1",

		OsnVersion: "1603",
		MapID:      "7",
		MapName:    "Peekaboo",
		MapTheme:   "2",

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
	err = osndb.Matches().Insert(db.NewMatchRecord(match))
	if err != nil {
		t.Errorf("error when inserting match metadata:\n%s", err)
	}

	log.Print(metadata, " => ", match)
	check_match(t, osndb, match)
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

func check_match(t *testing.T, db db.OsnDB, expected osn.LegacyMatch) {
	match, err := db.Matches().GetByName(string(expected.MatchHash))
	if err != nil {
		t.Errorf("error retrieving match %s\n%s\n", match.MatchHash, err)
	}

	if match.MatchIndex != expected.MatchIndex {
		t.Errorf("match index incorrect; got %d (expected %d)",
			match.MatchIndex, expected.MatchIndex)
	}
	if match.MatchHash != expected.MatchHash {
		t.Errorf("match ID incorrect; got %s (expected %s)",
			match.MatchHash, expected.MatchHash)
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

	if match.Version != expected.Version {
		t.Errorf("version %d != expected.version %d", match.Version, expected.Version)
	}
	if match.Status != expected.Status {
		t.Errorf("status %d != expected.status %d", match.Status, expected.Status)
	}

	if len(match.Players) != len(expected.Players) {
		t.Errorf("number of players %d != expected %d",
			len(match.Players), len(expected.Players))
		return
	}
	for i, player := range match.Players {
		if player.RowID != expected.Players[i].RowID {
			t.Errorf("player %d ID %d != expected %d", i, player.RowID, expected.Players[i].RowID)
		}
		if player.Name != expected.Players[i].Name {
			t.Errorf("player %d Name %s != expected %s", i, player.Name, expected.Players[i].Name)
		}
	}
}
