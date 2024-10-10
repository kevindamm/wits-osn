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
// github:kevindamm/wits-osn/cmd/fetch/backfill.go

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

// Backfill the contents of a TSV file into the Sqlite database instance.
func BackfillFromIndex(witsdb db.OsnDB, tsv_path string) error {
	reader, err := os.Open(tsv_path)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	values := strings.Split(scanner.Text(), "\t")

	columns := make([]string, 0)
	for _, col := range values {
		trimmed := strings.Trim(col, " ")
		columns = append(columns, trimmed)
	}
	indices := verify_columns(columns)
	log.Print(strings.Join(columns, ", "))

	linecount := 0
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		values = strings.Split(line, "\t")
		linecount += 1
		if len(values) != EXPECTED_TSV_COLUMN_COUNT {
			return fmt.Errorf("incorrect column count %d for record, line %d\n%s", len(values), linecount, line)
		}

		metadata := osn.LegacyReplayMetadata{
			GameID:      values[indices["game_id"]],
			NumPlayers:  num_players(values[indices["game_type"]]),
			LeagueMatch: league_match(values[indices["game_type"]]),
			Created:     values[indices["created"]],
			Season:      values[indices["season"]],
			OsnVersion:  values[indices["engine"]],
			MapID:       values[indices["map_id"]],
			MapName:     values[indices["map_name"]],
			TurnCount:   values[indices["turn_count"]],
		}

		metadata.MapID = values[indices["map_id"]]
		map_id := assert_uint8(values[indices["map_id"]])
		osnmap, err := witsdb.MapByID(uint8(map_id))
		if err != nil {
			return err
		}
		metadata.MapName = osnmap.Name

		playerIDs := split_list(values[indices["player_ids"]])
		playerNames := split_list(values[indices["player_names"]])
		playerLeagues := split_list(values[indices["player_leagues"]])
		playerRaces := split_list(values[indices["player_races"]])

		metadata.Player1_Name = playerNames[0]
		metadata.Player1_ID = playerIDs[0]
		metadata.Player1_League = playerLeagues[0]
		metadata.Player1_Race = playerRaces[0]

		metadata.Player2_Name = playerNames[1]
		metadata.Player2_ID = playerIDs[1]
		metadata.Player2_League = playerLeagues[1]
		metadata.Player2_Race = playerRaces[1]

		if metadata.NumPlayers == "4" {
			metadata.Player3_Name = playerNames[2]
			metadata.Player3_ID = playerIDs[2]
			metadata.Player3_League = playerLeagues[2]
			metadata.Player3_Race = playerRaces[2]

			metadata.Player4_Name = playerNames[3]
			metadata.Player4_ID = playerIDs[3]
			metadata.Player4_League = playerLeagues[3]
			metadata.Player4_Race = playerRaces[3]
		}

		log.Println(metadata)

		players := metadata.Players()
		for _, player := range players {
			record := db.PlayerRecord{Player: player}
			witsdb.Players().Insert(&record)
		}

		match := metadata.ToLegacyMatch()
		record := db.LegacyMatchRecord{LegacyMatch: match}
		witsdb.Matches().Insert(&record)

		if true {
			// TODO only processing one record until the schema transform is complete
			break
		}
	}

	return nil
}

func BackfillFromReplays(witsdb db.OsnDB, replays_path string) error {
	// TODO

	return nil

}

const EXPECTED_TSV_COLUMN_COUNT = 15

var EXPECTED_COLUMNS = map[string]bool{
	"game_id":        true,
	"game_type":      true,
	"season":         true,
	"created":        true,
	"player_names":   true,
	"player_ids":     true,
	"player_leagues": true,
	"player_races":   true,
	"map_id":         true,
	"map_name":       true,
	"turn_count":     true,
	"replay_fetched": true,
	"player_winner":  true,
	"engine":         true,
	"first_playerid": true,
}

func verify_columns(columns []string) map[string]int {
	if len(columns) != 15 {
		log.Fatalf("incorrect column count %d for TSV index", len(columns))
	}

	indices := make(map[string]int)
	for i, columnName := range columns {
		if expected := EXPECTED_COLUMNS[columnName]; expected {
			indices[columnName] = i
		} else {
			log.Fatalf("unrecognized column name %s in legacy index\n", columnName)
		}
	}
	if len(indices) != EXPECTED_TSV_COLUMN_COUNT {
		log.Fatalf("incorrect unique column count %d", len(indices))
	}

	return indices
}

// Convert string into an uint8 or fail with LOG(FATAL) << error.
func assert_uint8(intstr string) int8 {
	value64, err := strconv.ParseInt(intstr, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	if value64 < 0 || value64 > 255 {
		log.Fatal("expectued unsigned 8-bit integer, got " + intstr)
	}
	return int8(value64)
}

func num_players(gametype string) string {
	if gametype == "4" || gametype == "5" {
		return "4"
	} else if gametype != "0" {
		return "2"
	} else {
		return "0"
	}
}

func league_match(gametype string) string {
	gametypeint := assert_uint8(gametype)
	if gametypeint%2 > 0 {
		return "1"
	}
	return "0"
}

var reListCheck = regexp.MustCompile(`{(.*)(,.*)+}`)
var reListItem = regexp.MustCompile(`^("[^"]+"|[^"][^,}]*)[,}]`)

// Splits a string that is formatted as {...} wrapped comma-separated strings.
// Elements of the list that contain commas or spaces are wrapped in "-qoutes.
// The resulting list is two or four strings with any wrapping quotes removed.
// If there was an error in parsing, the empty []string{} list is returned.
//
// Reuses the backing string array for the slice of strings returned.
func split_list(liststr string) []string {
	list := make([]string, 0)
	if !reListCheck.Match([]byte(liststr)) {
		return []string{}
	}
	// Start with the first character after the open curly brace `{`.
	items := []byte(liststr[1:])

	// Parse items incrementally to more clearly avoid the quoted-commas.
	for reListItem.Match(items) {
		match := reListItem.FindIndex(items)
		item := string(items[:match[1]-1])
		list = append(list, item)
		// regexp match already includes the `,` or `}`
		items = items[match[1]:]
	}

	return list
}
