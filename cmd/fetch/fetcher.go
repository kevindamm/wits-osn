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
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	osn "github.com/kevindamm/wits-osn"
)

// Retrieves the data from an OSN index page and provides a channel for new
// (unique) match replay entries, or a nil channel and non-nil error.
func FetchIndexPage(pagenum int) ([]osn.LegacyMatch, error) {
	const URL = "http://osn.codepenguin.com/replays/getReplays/"

	values := url.Values{}
	if pagenum > 0 {
		values.Add("page", strconv.Itoa(pagenum))
		values.Add("ret_total", "false")
	}
	values.Add("limit", "19")
	values.Add("order", "created")
	values.Add("order_asc", "false")
	values.Add("list", "recent")

	response, err := http.PostForm(URL, values)
	if err != nil {
		return []osn.LegacyMatch{}, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return []osn.LegacyMatch{}, err
	}

	var index struct {
		// ignore other fields, we only care about the list of replays
		Total   string            `json:"total"`
		Replays []osn.LegacyMatch `json:"replays"`
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
func FetchReplay(game_id osn.GameID, filename string) error {
	const base_url = "http://osn.codepenguin.com/api/getReplay/"
	url := base_url + string(game_id)

	log.Print("Fetching ", url, " -> ", filename)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	wire_data, err := io.ReadAll(response.Body)
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

//	type IndexWithTotal struct {
//		Replays []Metadata `json:"replays"`
//		When    string     `json:"ts"`
//	}
//	responsebody, err := fetch_index(0)
//	if err != nil {
//		return
//	}
//
//	index := new(IndexWithTotal)
//	err = json.Unmarshal(responsebody, index)
//	if err != nil {
//		// unlikely error
//		return
//	}
//	total, err = strconv.Atoi(index.Total)
//	if err != nil {
//		// very unlikely error
//		return
//	}
//
//	return total, index.Replays
//}
