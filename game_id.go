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
// github:kevindamm/wits-go/osn/game_id.go

package osn

import (
	"encoding/json"
	"strings"
)

// Extension of a game's ID which automatically shortens it when marshaling.
type OsnGameID string

// Marshal the game ID (removing 48-character common prefix all games share).
func (id OsnGameID) MarshalJSON() ([]byte, error) {
	// Trim the common prefix off when writing the game ID.
	return json.Marshal(id.ShortID())
}

// Common prefix is overly redundant and does not need to keep taking up space.
const commonprefix string = "ahRzfm91dHdpdHRlcnNnYW1lLWhyZHIVCxIIR2FtZVJvb20Y"

func (id OsnGameID) ShortID() string {
	return strings.TrimPrefix(string(id), commonprefix)
}

func (id *OsnGameID) UnmarshalJSON(encoded []byte) error {
	// cut off infinite recursion
	var gameid string
	if err := json.Unmarshal(encoded, &gameid); err != nil {
		return err
	}
	*id = OsnGameID(gameid)
	return nil
}
