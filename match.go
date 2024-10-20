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
// github:kevindamm/wits-osn/match.go

package osn

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// The metadata of a single match between two or four players.
// Everything but the social signals (views/likes) and replay (player turns).
type LegacyMatch struct {
	MatchIndex  int64     `json:"-"`
	MatchHash   GameID    `json:"gameid"`
	OsnIndex    int       `json:"id,omitempty"`
	Competitive bool      `json:"competitive"`
	Season      int       `json:"season"`
	CreatedTime time.Time `json:"created"`
	MapID       int       `json:"mapid"`
	TurnCount   int       `json:"turn_count"`

	Version     int         `json:"engine"`
	FetchStatus FetchStatus `json:"-"`

	Players []PlayerRole `json:"players,omitempty"`
}

type GameOverData struct {
	Competitive Boolish           `json:"isLeagueMatch"`
	Online      Boolish           `json:"-"`
	Winners     []OsnPlayerUpdate `json:"winners"`
	Losers      []OsnPlayerUpdate `json:"losers"`
}

const UNKNOWN_MATCH_ID = GameID("")

var UNKNOWN_MATCH LegacyMatch = LegacyMatch{
	MatchHash: UNKNOWN_MATCH_ID,
}

// The (unaltered) representation of match-related metadata from OSN.
//
// This includes everything in the match entries of the json response
// for "recent replays".
type LegacyReplayMetadata struct {
	Index       string `json:"id"`            // integer index into sequential games
	GameID      string `json:"gameid"`        // hash of game creation
	NumPlayers  string `json:"gametype"`      // solo ("2") vs duo ("4")
	LeagueMatch string `json:"isleaguematch"` // "1" if true
	Created     string `json:"created"`       // 2012-08-05 14:33:21
	Season      string `json:"season"`        // integer ("1")

	// Runtime is parameterized by engine version and map definition.
	OsnVersion string `json:"engine"`     // integer #=< 1063, ("1000" is v1)
	MapID      string `json:"mapid"`      // relates to indexed maps, (id "4")
	MapName    string `json:"map_title"`  // display name ("Glitch"), redundant
	MapTheme   string `json:"map_raceid"` // enumeration, e.g. "1" is Feedback

	// These aren't persisted in the database but appear in the OSN index.
	TurnCount string `json:"turn_count"` // integer "34"
	ViewCount string `json:"viewcount"`  // integer "116"
	LikeCount string `json:"like_count"` // integer "0"

	Player1_ID     string `json:"p1_playerid"`   // integer "1",
	Player1_Name   string `json:"p1_playername"` // utf8 name "Syvan",
	Player1_League string `json:"p1_leagueid"`   // OsnLeagueEnum
	Player1_Race   string `json:"p1_raceid"`     // OsnRaceEnum
	Player1_Wins   string `json:"p1_winner"`     // boolish "1",
	Player1_BaseHP string `json:"p1_basehp"`     // #=< "5",

	Player2_ID     string `json:"p2_playerid"`   //: integer "2",
	Player2_Name   string `json:"p2_playername"` //: utf8 name "Alvendor",
	Player2_League string `json:"p2_leagueid"`   //: OsnLeagueEnum
	Player2_Race   string `json:"p2_raceid"`     //: OsnRaceEnum
	Player2_Wins   string `json:"p2_winner"`     //: boolish "0",
	Player2_BaseHP string `json:"p2_basehp"`     //: #=< "0",

	Player3_ID     string `json:"p3_playerid,omitempty"`   //: may be null
	Player3_Name   string `json:"p3_playername,omitempty"` //: may be null
	Player3_League string `json:"p3_leagueid,omitempty"`   //: may be null
	Player3_Race   string `json:"p3_raceid,omitempty"`     //: may be null
	Player3_Wins   string `json:"p3_winner,omitempty"`     //: may be null
	Player3_BaseHP string `json:"p3_basehp,omitempty"`     //: may be null

	Player4_ID     string `json:"p4_playerid,omitempty"`   //: may be null
	Player4_Name   string `json:"p4_playername,omitempty"` //: may be null
	Player4_League string `json:"p4_leagueid,omitempty"`   //: may be null
	Player4_Race   string `json:"p4_raceid,omitempty"`     //: may be null
	Player4_Wins   string `json:"p4_winner,omitempty"`     //: may be null
	Player4_BaseHP string `json:"p4_basehp,omitempty"`     //: may be null

	FirstPlayer string `json:"first_playerid"` // integer, matches playerid above
}

func (metadata *LegacyReplayMetadata) Players() []Player {
	numPlayers, err := strconv.Atoi(metadata.NumPlayers)
	if err != nil {
		return []Player{}
	}
	players := make([]Player, numPlayers)

	if numPlayers == 2 || numPlayers == 4 {
		playerRowID, err := strconv.ParseInt(metadata.Player1_ID, 10, 64)
		if err == nil {
			players[0] = NewPlayer(playerRowID, metadata.Player1_Name)
		} else {
			players[0] = UNKNOWN_PLAYER
		}

		playerRowID, err = strconv.ParseInt(metadata.Player2_ID, 10, 64)
		if err == nil {
			players[1] = NewPlayer(playerRowID, metadata.Player2_Name)
		} else {
			players[1] = UNKNOWN_PLAYER
		}
	}
	if numPlayers == 4 {
		playerRowID, err := strconv.ParseInt(metadata.Player3_ID, 10, 64)
		if err == nil {
			players[2] = NewPlayer(playerRowID, metadata.Player3_Name)
		} else {
			players[2] = UNKNOWN_PLAYER
		}

		playerRowID, err = strconv.ParseInt(metadata.Player4_ID, 10, 64)
		if err == nil {
			players[3] = NewPlayer(playerRowID, metadata.Player4_Name)
		} else {
			players[3] = UNKNOWN_PLAYER
		}
	}
	return players
}

// For go's [time.Parse] this must always be this same date and time.
const TimeLayout = "2006-01-02 15:04:05"

func (metadata *LegacyReplayMetadata) ToLegacyMatch() LegacyMatch {
	index := int64(0)
	if metadata.Index != "" {
		index = assert_int64(metadata.Index)
	}

	time, err := time.Parse(TimeLayout, metadata.Created)
	if err != nil {
		log.Fatalf("timestamp is not in expected format: %s", metadata.Created)
	}

	match := LegacyMatch{
		MatchIndex:  index,
		MatchHash:   GameID(metadata.GameID),
		OsnIndex:    int(index),
		Competitive: metadata.LeagueMatch == "1",
		Season:      int(assert_int64(metadata.Season)),
		StartTime:   time,
		MapID:       int(assert_int64(metadata.MapID)),
		TurnCount:   int(assert_int64(metadata.TurnCount)),
		Version:     int(assert_int64(metadata.OsnVersion)),
	}
	if metadata.NumPlayers == "4" {
		match.Players = make([]PlayerRole, 4)
	} else if metadata.NumPlayers == "2" {
		match.Players = make([]PlayerRole, 2)
	} else {
		log.Fatalf("metadata has unexpected NumPlayers %s", metadata.NumPlayers)
	}

	match.Players[0].RowID = assert_int64(metadata.Player1_ID)
	match.Players[0].Name = metadata.Player1_Name

	match.Players[1].RowID = assert_int64(metadata.Player2_ID)
	match.Players[1].Name = metadata.Player2_Name

	if metadata.NumPlayers == "4" {
		match.Players[2].RowID = assert_int64(metadata.Player3_ID)
		match.Players[2].Name = metadata.Player3_Name

		match.Players[3].RowID = assert_int64(metadata.Player4_ID)
		match.Players[3].Name = metadata.Player4_Name
	}

	return match
}

// Emits the JSON representation as well as the ID value which is typically not
// part of the JSON payload.  Useful for debugging, but the rowid `id` value is
// transient, it may change at the next VACUUM or repartitioning.
func (metadata LegacyReplayMetadata) String() string {
	if properties, err := json.Marshal(metadata); err == nil {
		return fmt.Sprintf("{ id: %s (%s),\n%s",
			metadata.Index, metadata.GameID,
			properties)
	} else {
		log.Printf("failed to marshal metadata (id: %s)\n%s", metadata.Index, err)
	}
	return metadata.Index
}

func assert_int64(image string) int64 {
	value, err := strconv.ParseInt(image, 10, 64)
	if err != nil {
		log.Fatalf("unexpected string value for int64: %s", image)
	}
	return value
}
