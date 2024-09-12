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
// github:kevindamm/wits-osn/cmd/fetch/index.go

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	osn "github.com/kevindamm/wits-osn"
)

type Metadata = osn.LegacyReplayMetadata

func FetchLatestIndex() (total int, replays []Metadata) {
	type IndexWithTotal struct {
		Total   string     `json:"total"`
		Replays []Metadata `json:"replays"`
		When    string     `json:"ts"`
	}
	responsebody, err := fetch_replaylist(0)
	if err != nil {
		return
	}

	index := new(IndexWithTotal)
	err = json.Unmarshal(responsebody, index)
	if err != nil {
		// unlikely error
		return
	}
	total, err = strconv.Atoi(index.Total)
	if err != nil {
		// very unlikely error
		return
	}

	return total, index.Replays
}

// Fetches the numbered page ; note that page 1 is the first 20 games
// and pages in the tens-of-thousands are more recent.
func FetchPage(page_number int) []Metadata {
	var data []byte
	var err error
	if page_number == 0 {
		// return (cached) latest page
		// TODO remove after dev
		data, err = os.ReadFile("../../testdata/index_latest.json")
		if err != nil {
			return []Metadata{}
		}
	} else {
		data, err = fetch_replaylist(page_number)
		if err != nil {
			return []Metadata{}
		}
	}

	var index struct {
		// ignore other fields, we only care about the list of replays
		Replays []Metadata `json:"replays"`
	}
	err = json.Unmarshal(data, &index)
	if err != nil {
		return []Metadata{}
	}

	return index.Replays
}

func fetch_replaylist(page_number int) ([]byte, error) {
	const URL = "http://osn.codepenguin.com/replays/getReplays/"

	values := url.Values{}
	if page_number > 0 {
		values.Add("page", strconv.Itoa(page_number))
	}
	values.Add("limit", "20")
	values.Add("order", "created")
	values.Add("order_asc", "false")
	values.Add("list", "recent")
	if page_number > 0 {
		values.Add("ret_total", "false")
	}

	response, err := http.PostForm(URL, values)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
