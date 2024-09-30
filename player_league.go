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
// github:kevindamm/wits-osn/player_league.go

package osn

import (
	"fmt"
	"log"
	"strconv"
)

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
	LeagueRange
)

func (league LeagueEnum) IsValid() bool {
	return league < LeagueRange
}

var league_names = []string{
	"UNKNOWN",
	"Fluffy",
	"Clever",
	"Gifted",
	"Master",
	"Supertitan",
}

func (league LeagueEnum) String() string {
	if !league.IsValid() {
		return "UNKNOWN"
	}
	return league_names[league]
}

// Parse the integer representation in the provided string, returns
// LEAGUE_UNKNOWN if there was an error or no league with that number.
func ParseLeague(str_uint string) LeagueEnum {
	value, err := strconv.Atoi(str_uint)
	if err != nil || !LeagueEnum(value).IsValid() {
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

func (rank LeagueRank) IsValid() bool { return true }
func (rank LeagueRank) String() string {
	return fmt.Sprintf("%d", rank)
}
