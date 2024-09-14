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
// github:kevindamm/wits-osn/cmd/fetch/fetcher.go

package main

import (
	osn "github.com/kevindamm/wits-osn"
)

// Retrieves JSON data from OSN URLs based on index contents of OSN listings.
// Maintains a persisted local state in a database (see [OsnWitsDB]) to reduce
// repeated retrievals.  Rate-limits retrievals of the index and match replays
// to reduce load on the OSN host.
type Fetcher interface {
	// Retrieves the data from an OSN index page and provides a channel for new
	// (unique) match replay entries, or a nil channel and non-nil error.
	FetchIndexPage(int) (<-chan osn.LegacyReplayMetadata, error)

	// Retrieves the match with indicated ID and saves it to a local file.
	// Returns the file path for the replay, or "" and a non-nil error.
	FetchReplay(string) (string, error)
}

// The implementation of Fetcher interface and its internal state.
type fetcher struct {
	db    OsnWitsDB
	path  string
	queue chan<- task
}

type task interface {
	Url() string
	Process([]byte) error
}

// Constructor function for a Fetcher instance.
func NewFetcher(witsdb OsnWitsDB, datapath string) Fetcher {
	taskchan := make(chan task)
	// TODO create cancelable goroutine for
	return &fetcher{witsdb, datapath, taskchan}
}

func (fetch *fetcher) FetchIndexPage(pagenum int) (<-chan osn.LegacyReplayMetadata, error) {
	// TODO
	return nil, nil
}

type IndexTask struct {
	url string
}

func (task IndexTask) Url() string { return task.url }
func (task IndexTask) Process(page []byte) error {
	// TODO
	return nil
}

/*
	type IndexWithTotal struct {
		Total   string     `json:"total"`
		Replays []Metadata `json:"replays"`
		When    string     `json:"ts"`
	}
	responsebody, err := fetch_index(0)
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
		data, err = fetch_index(page_number)
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
*/

func (fetch *fetcher) FetchReplay(string) (string, error) {
	// TODO
	return "", nil
}

type ReplayTask struct {
	url string
}

func (task ReplayTask) Url() string { return task.url }
func (task ReplayTask) Process(page []byte) error {
	// TODO
	return nil
}

/*
func fetch_index(page_number int) ([]byte, error) {
	const URL = "http://osn.codepenguin.com/replays/getReplays/"

	values := url.Values{}
	if page_number > -1 {
		values.Add("page", strconv.Itoa(page_number))
		values.Add("ret_total", "false")
	}
	values.Add("limit", "19")
	values.Add("order", "created")
	values.Add("order_asc", "false")
	values.Add("list", "recent")

	response, err := http.PostForm(URL, values)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
*/

/*
func fetch_replay(match_id string) ([]byte, error) {
	const URL = "http://osn.codepenguin.com/api/getReplay/"
	fmt.Printf("Fetching %s\n", URL+match_id)

	response, err := http.Get(URL + match_id)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
*/
