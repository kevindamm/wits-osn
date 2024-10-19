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
// github:kevindamm/wits-osn/db/maps_test.go

package db_test

import (
	"testing"

	"github.com/kevindamm/wits-osn/db"
)

func TestMapsTable(t *testing.T) {
	osndb := db.OpenOsnDB(":memory:")
	osndb.MustCreateAndPopulateTables()

	mapobj, err := osndb.MapByID(1)
	if err != nil {
		t.Errorf("could not find map ID 1: %s", err)
	} else if mapobj.MapID != 1 || mapobj.Name != "Machination" {
		t.Error("retrieved incorrect map for ID 1")
	}

	mapobj, err = osndb.MapByName("foundry")
	if err != nil {
		t.Errorf("could not find map ID 3 (Foundry): %s", err)
	}
	if mapobj.MapID != 3 || mapobj.Name != "Foundry" {
		t.Error("retrieved incorrect map for ID 3")
	}

	// Maps are read-only and no INSERT/DELETE interface is exposed.
}
