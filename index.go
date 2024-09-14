// Copyright (c) 2024 Kevin Damm
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	GameType    string `json:"gametype"`      // solo "2" vs duo "4"
	LeagueMatch string `json:"isleaguematch"` // "1" if true
	Created     string `json:"created"`       // 2012-08-05 14:33:21
	Season      string `json:"season"`        // integer "1"

	// Runtime is parameterized by engine version and map definition.
	WitsVersion string `json:"engine"`     // integer #=< 1063, "1000" is v1
	MapID       string `json:"mapid"`      // indexed, e.g. "4"
	MapName     string `json:"map_title"`  // display name, e.g. "Glitch"
	MapTheme    string `json:"map_raceid"` // enumeration, e.g. "1"

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
	numPlayers, err := strconv.Atoi(metadata.GameType)
	if err != nil {
		return []Player{}
	}
	if numPlayers%2 == 1 {
		numPlayers -= 1
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
			players[3] = NewPlayer(playerRowID, metadata.Player3_Name)
		} else {
			players[3] = UNKNOWN_PLAYER
		}

		playerRowID, err = strconv.Atoi(metadata.Player4_ID)
		if err == nil {
			players[4] = NewPlayer(playerRowID, metadata.Player4_Name)
		} else {
			players[4] = UNKNOWN_PLAYER
		}
	}
	return players
}

func (metadata *LegacyReplayMetadata) ToLegacyMatch() LegacyMatch {
	match := LegacyMatch{}

	return match
}
