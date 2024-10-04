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
// github:kevindamm/wits-osn/player_standing.go

package osn

import "fmt"

// An ELO-like measurement [points], and the player's league + standings status.
type PlayerStanding struct {
	League LeagueEnum
	Rank   LeagueRank
	Points uint16
	Delta  int8 // difference since [points] of previous standings.
}

func NewStanding(league LeagueEnum, rank LeagueRank, points uint16, delta int8) (PlayerStanding, error) {
	if !league.IsValid() {
		return UnknownStanding(), fmt.Errorf("invalid league value %d", league)
	}
	if !rank.IsValid() {
		return UnknownStanding(), fmt.Errorf("invalid rank value %d", rank)
	}
	return PlayerStanding{league, rank, points, delta}, nil
}

func (ranked PlayerStanding) PointsBefore() uint16 {
	return uint16(int(ranked.Points) - int(ranked.Delta))
}

func UnknownStanding() PlayerStanding {
	return PlayerStanding{}
}
