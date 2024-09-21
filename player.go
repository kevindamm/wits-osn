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

import (
	"fmt"
	"log"
	"strconv"
)

type Player struct {
	ID   PlayerID
	Name string
}

// Simple (no GCID) constructor for a Player instance.
func NewPlayer(id int, name string) Player {
	return Player{ID: PlayerID{RowID: id}, Name: name}
}

// Represents an assignment of a Player (identifier) with match participation.
type PlayerRole struct {
	Player
	UnitRace    UnitRaceEnum
	PlayerColor PlayerColorEnum

	AP        int
	BaseTheme int

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
type PlayerStanding interface {
	League() LeagueEnum
	Rank() LeagueRank

	PointsBefore() uint16
	PointsAfter() uint16
	Delta() int8
}

func NewStanding(league LeagueEnum, rank LeagueRank, points uint16, delta int8) (PlayerStanding, error) {
	if !league.Valid() {
		return nil, fmt.Errorf("invalid league value %d", rank)
	}
	if rank < 0 || rank >= 128 {
		return nil, fmt.Errorf("invalid standings rank value %d", rank)
	}
	return standing{league, rank, points, delta}, nil
}

type standing struct {
	league LeagueEnum
	rank   LeagueRank
	points uint16
	delta  int8
}

func (ranked standing) League() LeagueEnum  { return ranked.league }
func (ranked standing) Rank() LeagueRank    { return ranked.rank }
func (ranked standing) PointsAfter() uint16 { return ranked.points }
func (ranked standing) Delta() int8         { return ranked.delta }
func (ranked standing) PointsBefore() uint16 {
	return uint16(int(ranked.points) - int(ranked.delta))
}

// FROZEN
type LeagueEnum uint8

// Using iota here because I don't plan on ever changing the OSN representation.
// See github.com/kevindamm/wits-go for a more forward-compatible enumeration.
const (
	LEAGUE_UNKNOWN LeagueEnum = iota
	LEAGUE_FLUFFY
	LEAGUE_CLEVER
	LEAGUE_GIFTED
	LEAGUE_MASTER
	LEAGUE_SUPERTITAN
)

func (league LeagueEnum) Valid() bool {
	return uint8(league) > 0 && uint8(league) <= uint8(LEAGUE_SUPERTITAN)
}

// Parse the integer representation in the provided string, returns
// LEAGUE_UNKNOWN if there was an error or no league with that number.
func ParseLeague(str_uint string) LeagueEnum {
	value, err := strconv.Atoi(str_uint)
	if err != nil || !LeagueEnum(value).Valid() {
		log.Printf("unrecognized league %s\n%s", str_uint, err)
		return LEAGUE_UNKNOWN
	}
	return LeagueEnum(value)
}

// The player's rank within their current group.
// Players are placed in groups of around 100 when entering a league.
// Each of these divisions is given a name but historical data of that is
// not included in the replay or index data sent from OSN, and is not relevant
// enough to the replay analysis to be archived with the replay data.
type LeagueRank uint8

// This is actually arbitrary but provided for internal consistency.
type PlayerColorEnum int

const (
	COLOR_BLUE  PlayerColorEnum = 1
	COLOR_RED   PlayerColorEnum = 2
	COLOR_GREEN PlayerColorEnum = 3
	COLOR_GOLD  PlayerColorEnum = 4
)

func (color PlayerColorEnum) String() string {
	return []string{
		"BLUE", "RED", "GREEN", "GOLD",
	}[int(color)]
}
