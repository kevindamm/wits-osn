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
// github:kevindamm/wits-osn/db/matches_test.go

package db_test

import (
	"log"
	"testing"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

func TestMatchesTable(t *testing.T) {
	osndb := db.OpenOsnDB(":memory:")
	osndb.MustCreateAndPopulateTables()

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
	err := osndb.Matches().Insert(db.MakeMatchRecord(match))
	if err != nil {
		t.Errorf("error when inserting match metadata:\n%s", err)
	}

	log.Print(metadata, " => ", match)
	check_match(t, osndb, match)
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
	if match.FetchStatus != expected.FetchStatus {
		t.Errorf("status %d != expected.status %d", match.FetchStatus, expected.FetchStatus)
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
