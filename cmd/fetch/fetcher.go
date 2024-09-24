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
// github:kevindamm/wits-osn/fetcher.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	osn "github.com/kevindamm/wits-osn"
	"github.com/kevindamm/wits-osn/db"
)

type Fetcher interface {
	FetchNewReplayIDs(db.OsnDB) (<-chan osn.GameID, <-chan error)
	FetchReplay(osn.GameID, string) error
}

func NewFetcher(waitSeconds uint) Fetcher {
	// <3 be kind to your hosts <3 <3 <3
	if waitSeconds < 3 {
		waitSeconds = 3
	}
	return &fetcher{time.Now(), waitSeconds, fetch_index, fetch_replay}
}

type fetcher struct {
	// Rate-limiting of fetches across resource types.
	fetchedAt   time.Time
	waitSeconds uint

	// Default to making a network request, may be mocked out by tests.
	fetch_index  func(pagenum int) ([]byte, error)
	fetch_replay func(pageurl string) ([]byte, error)
}

func (fetcher *fetcher) wait() {

}

func (fetcher *fetcher) FetchNewReplayIDs(db db.OsnDB) (<-chan osn.GameID, <-chan error) {
	errchan := make(chan error)
	idchan := make(chan osn.GameID)

	go func() {
		defer close(idchan)
		defer close(errchan)

		i := 0
		replays := []osn.LegacyMatch{}
		for !all_fetched(replays, db) {
			matches, err := fetcher.fetch_and_parse_index(i)
			if err != nil {
				errchan <- err
				break
			}
			for _, match := range matches {
				// TODO write to database if it doesn't already exist

				// TODO don't emit to channel if it had already been fetched
				// (may be from page shift or from already-retrieved history)
				idchan <- osn.GameID(match.MatchHash)
			}

			i += 1
		}
	}()

	return idchan, errchan
}

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
const ReplayFormURL = "http://osn.codepenguin.com/replays/getReplays/"

func all_fetched(matches []osn.LegacyMatch, db db.OsnDB) bool {
	// TODO
	return true
}

// Retrieves the data from an OSN index page and provides a channel for new
// (unique) match replay entries, or a nil channel and non-nil error.
func fetch_index(pagenum int) ([]byte, error) {
	// Maintain the same ordering and representation as the browser interface.
	values := url.Values{}
	if pagenum > 0 {
		values.Add("page", strconv.Itoa(pagenum))
	}
	values.Add("limit", "20")
	values.Add("order", "created")
	values.Add("order_asc", "false")
	values.Add("list", "recent")
	if pagenum > 0 {
		values.Add("ret_total", "false")
	}

	request, err := http.NewRequest("POST", ReplayFormURL, strings.NewReader(values.Encode()))
	if err != nil {
		return []byte{}, err
	}
	request.Header.Set("User-Agent", UserAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func (fetcher *fetcher) fetch_and_parse_index(pagenum int) ([]osn.LegacyMatch, error) {
	data, err := fetcher.fetch_index(pagenum)
	if err != nil {
		return []osn.LegacyMatch{}, err
	}

	var index struct {
		Total   string            `json:"total,omitempty"`
		Replays []osn.LegacyMatch `json:"replays"`
		// Ignore server timestamp, it isn't of any importance (and it drifts).
		//When    string            `json:"ts"`
	}
	err = json.Unmarshal(data, &index)
	if err != nil {
		return []osn.LegacyMatch{}, err
	}

	return index.Replays, nil
}

// Retrieves the match with indicated ID and saves it to a local file.
// Returns the game ID for the replay, or "" and a non-nil error.
func (fetcher *fetcher) FetchReplay(game_id osn.GameID, filename string) error {
	url := fmt.Sprintf("http://osn.codepenguin.com/api/getReplay/%s", game_id)

	log.Print("Fetching ", url, " -> ", filename)
	wire_data, err := fetcher.fetch_replay(url)
	if err != nil {
		return err
	}

	_, encoded, err := osn.ParseRawReplay(wire_data)
	if err != nil {
		return err
	}

	os.WriteFile(filename, encoded, 0644)
	return nil
}

func fetch_replay(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
