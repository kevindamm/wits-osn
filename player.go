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
// github:kevindamm/wits-osn/player.go

package osn

type Player struct {
	ID   PlayerID
	Name string
}

type PlayerRole struct {
	Player
	TurnOrder  int
	RankBefore PlayerStanding
	RankAfter  PlayerStanding
}

// Represents both (or either of) the internal identity and the player's GCID.
// A zero value for either one will indicate missing or unknown value.
type PlayerID struct {
	RowID int
	GCID  string
}

var UNKNOWN_PLAYER Player = Player{
	ID:   PlayerID{RowID: 0, GCID: ""},
	Name: "",
}

// An ELO-like measurement and a current league + standings.
type PlayerStanding struct {
	League LeagueEnum
	Rank   LeagueRank

	Points uint16
	Delta  int8
}

type LeagueEnum uint8

const (
	LEAGUE_UNKNOWN    LeagueEnum = 0
	LEAGUE_FLUFFY     LeagueEnum = 1
	LEAGUE_CLEVER     LeagueEnum = 2
	LEAGUE_GIFTED     LeagueEnum = 3
	LEAGUE_MASTER     LeagueEnum = 4
	LEAGUE_SUPERTITAN LeagueEnum = 5
)

// The player's rank within their current group.
//
type LeagueRank uint8
