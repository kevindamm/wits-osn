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
// github:kevindamm/wits-osn/units.go

package osn

type UnitStatus struct {
	AltHealth uint            `json:"altHealth"`
	Class     UnitClass       `json:"class"`
	Color     PlayerColorEnum `json:"color"`

	HasAttacked    Boolish `json:"hasAttacked"`
	HasMoved       Boolish `json:"hasMoved"`
	HasTransformed Boolish `json:"hasTransformed"`

	Health     uint    `json:"health"`
	Identifier uint    `json:"identifier"`
	IsAlt      Boolish `json:"isAlt"`
	Owner      uint    `json:"owner"` // player turn order?
	Team       uint    `json:"team"`

	Parent      int `json:"parent"`      // -1 if no parentage
	SpawnedFrom int `json:"spawnedFrom"` // -1 if from spawn tile

	PositionI uint         `json:"positionI"`
	PositionJ uint         `json:"positionJ"`
	UnitRace  UnitRaceEnum `json:"race"`
}

// TODO populate enum
type UnitClass int
