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
// github:kevindamm/wits-osn/fetch_status.go

package osn

// There is a finite, unchanging set of status values
// based on the progress of processing a sing replay
// from reading its listing through transforming its
// encoding and reducing combinatoric symmetries.
type FetchStatus uint8

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
	StatusRange
)

func (status FetchStatus) IsValid() bool {
	return uint8(status) <= uint8(STATUS_LEGACY)
}

var status_names = []string{
	"UNKNOWN",
	"LISTED",
	"FETCHED",
	"UNWRAPPED",
	"CONVERTED",
	"CANONICAL",
	"VALIDATED",
	"INDEXED",
	"INVALID",
	"LEGACY",
}

func (status FetchStatus) String() string {
	if !status.IsValid() {
		return "UNKNOWN"
	}
	return status_names[status]
}
