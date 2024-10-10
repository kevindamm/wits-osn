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
	"strings"

	osn "github.com/kevindamm/wits-osn"
	_ "github.com/mattn/go-sqlite3"
)

// Interface for simplifying the interaction with a backing database.
//
// DB includes match metadata, map identities, player history and replay index.
type OsnDB interface {
	MustCreateAndPopulateTables() // Create tables or die trying.
	Close()                       // Closes the database and attached resources.

	MapByID(id uint8) (osn.LegacyMap, error)
	MapByName(name string) (osn.LegacyMap, error)

	Players() MutableTable[*PlayerRecord]
	Matches() MutableTable[*LegacyMatchRecord]
	Standings() MutableTable[*StandingsRecord]

	UpdateMatchStatus(osn.GameID, osn.FetchStatus) error
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

	status  EnumTable[osn.FetchStatus]
	leagues EnumTable[osn.LeagueEnum]
	races   EnumTable[osn.UnitRaceEnum]

	maps Table[*LegacyMapRecord]

	players   MutableTable[*PlayerRecord]
	matches   MutableTable[*LegacyMatchRecord]
	roles     MutableTable[*PlayerRoleRecord]
	standings MutableTable[*StandingsRecord]
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
	// TODO: Perhaps track all open requests via [Context] and cancel them too.
	db.sqldb.Close()
}

func (db *osndb) MapByID(id uint8) (osn.LegacyMap, error) {
	record, err := db.maps.Get(int64(id))
	if err != nil {
		return osn.UnknownMap(), err
	}
	return osn.LegacyMap(*record), nil
}

func (db *osndb) MapByName(name string) (osn.LegacyMap, error) {
	record, err := db.maps.GetByName(name)
	if err != nil {
		return osn.UnknownMap(), err
	}
	return osn.LegacyMap(*record), nil
}

func (db *osndb) Players() MutableTable[*PlayerRecord]      { return db.players }
func (db *osndb) Matches() MutableTable[*LegacyMatchRecord] { return db.matches }
func (db *osndb) Standings() MutableTable[*StandingsRecord] { return db.standings }

func (db *osndb) UpdateMatchStatus(matchID osn.GameID, status osn.FetchStatus) error {
	// TODO
	return nil
}

// Only needs to be called once at database setup.  Also closes the database.
// Will LOG(FATAL) an error if creation or initialization fail, with SQL error.
func (db *osndb) MustCreateAndPopulateTables() {
	for _, table := range []TableSql{
		db.status,
		db.leagues,
		db.races,
		db.maps,
		db.players,
		db.matches,
		db.roles,
		db.standings,
	} {
		createsql := table.SqlCreate()
		log.Println(createsql)
		must_execsql(db.sqldb, createsql)

		initsql := table.SqlInit()
		if initsql != "" {
			log.Println(initsql)
			statements := strings.Split(initsql, ";")
			for _, sql := range statements {
				must_execsql(db.sqldb, sql)
			}
		}
	}
}
