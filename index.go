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
// github:kevindamm/wits-osn/replay.go

package osn

import "strconv"

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
	WitsVersion string `json:"engine"`     // integer #=< 1063, ("1000" is v1)
	MapID       string `json:"mapid"`      // relates to indexed maps, (id "4")
	MapName     string `json:"map_title"`  // display name ("Glitch"), redundant
	MapTheme    string `json:"map_raceid"` // enumeration, e.g. "1" is Feedback

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
		playerRowID, err := strconv.Atoi(metadata.Player1_ID)
		if err == nil {
			players[0] = NewPlayer(playerRowID, metadata.Player1_Name)
		} else {
			players[0] = UNKNOWN_PLAYER
		}

		playerRowID, err = strconv.Atoi(metadata.Player2_ID)
		if err == nil {
			players[1] = NewPlayer(playerRowID, metadata.Player2_Name)
		} else {
			players[1] = UNKNOWN_PLAYER
		}
	}
	if numPlayers == 4 {
		playerRowID, err := strconv.Atoi(metadata.Player3_ID)
		if err == nil {
			players[2] = NewPlayer(playerRowID, metadata.Player3_Name)
		} else {
			players[2] = UNKNOWN_PLAYER
		}

		playerRowID, err = strconv.Atoi(metadata.Player4_ID)
		if err == nil {
			players[3] = NewPlayer(playerRowID, metadata.Player4_Name)
		} else {
			players[3] = UNKNOWN_PLAYER
		}
	}
	return players
}

func (metadata *LegacyReplayMetadata) ToLegacyMatch() LegacyMatch {
	match := LegacyMatch{}

	return match
}
