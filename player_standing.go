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

// An ELO-like measurement and a current league + standings.
type PlayerStanding interface {
	League() LeagueEnum
	Rank() LeagueRank

	PointsBefore() uint16
	PointsAfter() uint16
	Delta() int8
}

func NewStanding(league LeagueEnum, rank LeagueRank, points uint16, delta int8) (PlayerStanding, error) {
	if !league.IsValid() {
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
