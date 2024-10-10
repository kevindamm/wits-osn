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
// github:kevindamm/wits-osn/cmd/fetch/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

func main() {
	db_path := flag.String("db-path", ".data/osn.db",
		"path where the sqlite3 database will be written")
	create_tables := flag.Bool("create-tables", false,
		"create the table schemata before running, including enum values")
	backfill_tsv := flag.String("backfill-tsv", "",
		"path to a TSV file containing a legacy backup of replay file metadata")
	backfill_replays := flag.String("backfill-replays", "",
		"path to the parent directory where JSON files of wits replays are found")
	out_path := flag.String("out", ".data/replays/",
		"path where JSON formatted replay data is written to")

	flag.Parse()

	witsdb := db.OpenOsnDB(*db_path)
	defer witsdb.Close()

	if *create_tables {
		log.Println("creating DB tables...")
		witsdb.MustCreateAndPopulateTables()
	}

	if len(*backfill_tsv) > 0 {
		log.Println("back-filling from legacy DB...")
		assert_nilerr(BackfillFromIndex(witsdb, *backfill_tsv))
	}
	if len(*backfill_replays) > 0 {
		log.Println("back-filling from legacy replays...")
		assert_nilerr(BackfillFromReplays(witsdb, *backfill_replays))
	}

	assert_nilerr(os.MkdirAll(*out_path, 0755))

	// Fetch listing of recent (unaccounted-for) replays
	fetcher := NewFetcher(5)
	replay_index, errs := fetcher.FetchNewReplayIDs(witsdb)
	count := 0

	// Fetch replays that haven't been fetched already
	for {
		select {
		case replayID, ok := <-replay_index:
			if ok {
				replay_path := path.Join(*out_path,
					fmt.Sprintf("%s.json", replayID.ShortID()))
				fmt.Printf("fetching %s -> %s", replayID.ShortID(), replay_path)
				err := fetcher.FetchReplay(replayID, replay_path)
				if err == nil {
					count += 1
					witsdb.UpdateMatchStatus(replayID, osn.STATUS_FETCHED)
				} else {
					fmt.Println("ERROR: ", err)
				}
			} else {
				replay_index = nil
			}
		case err, ok := <-errs:
			if ok {
				fmt.Println("ERROR while fetching index pages")
				fmt.Println(err)
				return
			} else {
				errs = nil
			}
		}
		if replay_index == nil && errs == nil {
			break
		}
	}

	fmt.Println("fetch of recent replays completed")
	fmt.Printf("%d new replays fetched", count)
}

func assert_nilerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
