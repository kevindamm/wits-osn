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
	"strconv"
	"strings"

	osn "github.com/kevindamm/wits-osn"
)

// Backfill the contents of a TSV file into the Sqlite database instance.
func BackfillFromIndex(witsdb WitsDB, tsv_path string) error {
	reader, err := os.Open(tsv_path)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	values := strings.Split(scanner.Text(), ", ")

	columns := make([]string, 0)
	for _, col := range values {
		trimmed := strings.Trim(col, " ")
		columns = append(columns, trimmed)
	}
	indices := verify_columns(columns)

	linecount := 0
	for scanner.Scan() {
		line := scanner.Text()
		values = strings.Split(line, "\t")
		linecount += 1
		if len(values) != EXPECTED_TSV_COLUMN_COUNT {
			if linecount >= 1831407 {
				break
			}
			return fmt.Errorf("incorrect column count %d for record, line %d\n%s", len(values), linecount, line)
		}

		metadata := osn.LegacyReplayMetadata{

			GameID: values[indices["game_id"]],
			//GameType:
		}
		curmap, err := witsdb.Map(assert_int(values[indices["map_id"]]))
		if err != nil {
			return err
		}
		metadata.MapID = strconv.Itoa(curmap.MapID)

		//players := metadata.Players()

		match := metadata.ToLegacyMatch()
		witsdb.InsertMatch(match)

		if true {
			// TODO only processing one record until the schema transform is complete
			break
		}
	}

	return nil
}

func BackfillFromReplays(witsdb WitsDB, replays_path string) error {
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
	"player_orders":  true,
	"map_id":         true,
	"map_name":       true,
	"turn_count":     true,
	"replay_fetched": true,
	"player_winner":  true,
	"engine":         true,
}

func verify_columns(columns []string) map[string]int {
	if len(columns) != 15 {
		log.Fatalf("incorrect column count %d for TSV index", len(columns))
	}

	indices := make(map[string]int)
	for i, columnName := range columns {
		if _, ok := EXPECTED_COLUMNS[columnName]; ok {
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

// Convert string into an integer or fail with LOG(FATAL) << error.
func assert_int(intstr string) int {
	value, err := strconv.Atoi(intstr)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
