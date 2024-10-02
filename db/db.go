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
// github:kevindamm/wits-osn/db/db.go

package db

import (
	"database/sql"
	"log"

	osn "github.com/kevindamm/wits-osn"
	_ "github.com/mattn/go-sqlite3"
)

// Interface for simplifying the interaction with a backing database.
//
// DB includes match metadata, map identities, player history and replay index.
type OsnDB interface {
	Close() // Closes the database, including open rows and prepared statements.

	Map(name string) (osn.LegacyMap, error)

	Players() Table[*PlayerRecord]
	Matches() Table[*LegacyMatchRecord]
	Standings() Table[*StandingsRecord]
}

// Opens a Sqlite db at indicated path and prepares queries.
// Does not create tables, only prepares the connection and statements.
//
// Asserts db exists and basic queries can be constructed.
// LOG(FATAL) on any error, database remains unmodified.
func OpenOsnDB(filepath string) OsnDB {
	osndb, err := open_database(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return OsnDB(osndb)
}

type osndb struct {
	sqldb *sql.DB

	// read-only tables
	status  EnumTable[osn.FetchStatus]
	leagues EnumTable[osn.LeagueEnum]
	races   EnumTable[osn.UnitRaceEnum]

	maps Table[LegacyMapRecord]

	// writable tables
	players   Table[*PlayerRecord]
	matches   Table[*LegacyMatchRecord]
	roles     Table[*PlayerRoleRecord]
	standings Table[*StandingsRecord]
}

// Opens a connection to the database but does not prepare any queries.
//
// Useful for getting a connection to create the required tables, Callers
// should prefer using [OpenOsnDB()], [CreateTables()], or methods on [OsnDB].
//
// Internally, this is used to get a connection without preparing queries (which
// would fail when required tables had not been created yet).
func open_database(filepath string) (*osndb, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	osndb := new(osndb)
	osndb.sqldb = db

	// We can initialize the table mappings without preparing queries.
	osndb.status = MakeEnumTable("fetch_status",
		osn.EnumValuesFor(osn.FetchStatusRange))
	osndb.leagues = MakeEnumTable("player_leagues",
		osn.EnumValuesFor(osn.LeagueRange))
	osndb.races = MakeEnumTable("races",
		osn.EnumValuesFor(osn.UnitRaceRange))

	osndb.maps = MakeMapsTable(osndb.sqldb)
	osndb.players = MakePlayersTable(osndb.sqldb)
	osndb.matches = MakeMatchesTable(osndb.sqldb)
	osndb.roles = MakeRolesTable(osndb.sqldb)
	osndb.standings = MakeStandingsTable(osndb.sqldb)

	return osndb, nil
}

func (db *osndb) Close() {
	// TODO close open statements & rows, cancel any pending contexts
	db.sqldb.Close()
}

func (db *osndb) Map(name string) (osn.LegacyMap, error) {
	// TODO
	return osn.UnknownMap(), nil
}

func (db *osndb) Players() Table[*PlayerRecord]      { return db.players }
func (db *osndb) Matches() Table[*LegacyMatchRecord] { return db.matches }
func (db *osndb) Standings() Table[*StandingsRecord] { return db.standings }
