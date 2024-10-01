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

// Represents an assignment of a Player (identifier) with match participation.
type PlayerRole struct {
	Player    `json:"player" orm:"fk(players)"`
	UnitRace  UnitRaceEnum    `json:"race" orm:"fk(races)"`
	BaseTheme int             `json:"theme" orm:"-"`
	TurnOrder PlayerColorEnum `json:"color" orm:"turn_order"`

	// The remaining actions for the player at the latest turn.
	//
	// This is useful for maintaining a view of the state derived from all of the
	// player's previous turns but doesn't need to be persisted in the database.
	Actions uint `json:"wits" orm:"-"`

	// These are derived from the standings table which refer to the role record.
	RankBefore PlayerStanding `json:"rank_prev" orm:"--from fk(standings.until)"`
	RankAfter  PlayerStanding `json:"rank_next,omitempty" orm:"--from fk(standings.after)"`
}
