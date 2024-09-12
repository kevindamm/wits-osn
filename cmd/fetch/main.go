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
	"log"
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

	flag.Parse()

	//replays := FetchPage(0)
	//for _, replay := range replays {
	//	replayjson, err := json.Marshal(replay)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	fmt.Println(string(replayjson))
	//}

	//replay_example, err := fetch_replay("ag5vdXR3aXR0ZXJzZ2FtZXIQCxIIR2FtZVJvb20Y1pBfDA")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(string(replay_example))
	//}
	if *create_tables {
		log.Println("creating DB tables...")
		CreateTablesAndClose(*db_path)
	}

	witsdb := OpenWitsDB(*db_path)
	defer witsdb.Close()

	if len(*backfill_tsv) > 0 {
		log.Println("back-filling from legacy DB...")
		assertNil(BackfillFromIndex(witsdb, *backfill_tsv))
	}
	if len(*backfill_replays) > 0 {
		log.Println("back-filling from legacy replays...")
		assertNil(BackfillFromReplays(witsdb, *backfill_replays))
	}

	// TODO fetch listing of recent (unaccounted-for) replays

	// TODO fetch replays that haven't been fetched already

}
