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
// github:kevindamm/wits-osn

package osn

import "encoding/json"

// This is the format as returned by the web service for a single game replay.
//
// It is a shallow wrapper around the actual game replay, containing four
// dictionary entries, one of which is a string representation of the game.
type WireFormat struct {
	Wrapper outer `json:"viewResponse"` // Path reference of an `.hxm` file.
}

type outer struct {
	Wrapper string `json:"gameState"` // Wrapper around the turn's frames.
	Found   bool   `json:"foundRoom"` // Always true - a game lobby was found.
	RoomID  string `json:"room"`      // The game (and game replay) identifier.
}

type inner struct {
	Wrapper string `json:"gameState"`
}

func ParseReplay(filedata []byte) (string, []byte, error) {
	var on_wire WireFormat
	var gamestate inner
	var replay LegacyMatchReplay
	err := json.Unmarshal(filedata, &on_wire)
	if err != nil {
		return "", []byte{}, err
	}
	err = json.Unmarshal([]byte(on_wire.Wrapper.Wrapper), &gamestate)
	if err != nil {
		return on_wire.Wrapper.RoomID, []byte{}, err
	}
	err = json.Unmarshal([]byte(gamestate.Wrapper), &replay)
	if err != nil {

	}
	bytes, err := json.Marshal(replay)
	return on_wire.Wrapper.RoomID, bytes, err
}

// The slightly-flattened complete representation of a match replay from OSN.
type LegacyMatchReplay struct {
	OsnMatchID string `json:"game_id"`
}
