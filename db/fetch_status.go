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
// github:kevindamm/wits-osn/db/fetch_status.go

package db

type FetchStatus int

const (
	STATUS_UNKNOWN FetchStatus = iota
	STATUS_LISTED
	STATUS_FETCHED
	STATUS_UNWRAPPED
	STATUS_CONVERTED
	STATUS_CANONICAL
	STATUS_VALIDATED
	STATUS_INDEXED
	STATUS_INVALID
	STATUS_LEGACY
)

//		`CREATE TABLE "fetch_status" (
//      "id"    INTEGER PRIMARY KEY,
//      "name"  VARCHAR(10) NOT NULL
//    ) WITHOUT ROWID;`,
//
//		`INSERT INTO fetch_status VALUES
//      (0, "UNKNOWN"),
//      (1, "LISTED"),
//      (2, "FETCHED"),
//      (3, "UNWRAPPED"),
//      (4, "CONVERTED"),
//      (5, "CANONICAL"),
//      (6, "VALIDATED"),
//      (7, "INDEXED"),
//      (8, "INVALID"),
//      (9, "LEGACY");`,
