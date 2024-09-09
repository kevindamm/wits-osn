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

// The (unaltered) representation of match-related metadata from OSN.
type LegacyReplayMetadata struct {
	Index       string `json:"id"`            // integer index into sequential games
	GameID      string `json:"gameid"`        // hash of game creation
	GameType    string `json:"gametype"`      // solo "2" vs duo "4"
	LeagueMatch string `json:"isleaguematch"` // "1" if true
	Created     string `json:"created"`       // 2012-08-05 14:33:21

	MapID   string `json:"mapid"`      // indexed, e.g. "4"
	MapName string `json:"map_title"`  // display name, e.g. "Glitch"
	RaceID  string `json:"map_raceid"` // enumeration, e.g. "1"

	TurnCount   string `json:"turn_count"` // integer "34"
	ViewCount   string `json:"viewcount"`  // integer "116"
	LikeCount   string `json:"like_count"` // integer "0"
	WitsVersion string `json:"engine"`     // integer #=< 1063, "1000" is v1
	WitsSeason  string `json:"season"`     // integer "1"

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

	FirstPlayer string `json:"first_playerid"` // "2" matches playerid above
}

// The slightly-flattened complete representation of a match replay from OSN.
type LegacyMatchReplay struct {
	OsnGameID string `json:"game_id"`
	MapName   string `json:"map_name"`
	MapTheme  string `json:"map_theme"`
}
