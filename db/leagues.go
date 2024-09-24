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
// github:kevindamm/wits-osn/db/leagues.go

package db

//		// Enumerative relation for the different ranked leagues.
//		// Players may be promoted or demoted based on win/loss records.
//		`CREATE TABLE "leagues" (
//      "league_id"    INTEGER PRIMARY KEY,
//      "league_name"  TEXT NOT NULL,
//
//      UNIQUE (league_name) ON CONFLICT IGNORE
//    ) WITHOUT ROWID;`,
//
//		`INSERT INTO leagues VALUES
//      (0, "UNKNOWN"),
//      (1, "Fluffy"),
//      (2, "Clever"),
//      (3, "Gifted"),
//      (4, "Master"),
//      (5, "SuperTitan");`,
//
