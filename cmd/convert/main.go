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
// github.com:kevindamm/wits-osn/cmd/convert/main.go

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	// truncate database backup to contain only the metadata
	reader, err := os.Open(".data/osnwits_db.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	values := strings.Split(scanner.Text(), ", ")

	columns := make([]string, 0)
	for _, col := range values {
		trimmed := strings.Trim(col, " ")
		columns = append(columns, trimmed)
	}
	column_count := len(columns)

	writer, err := os.Create(".data/osn_index.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	writer.WriteString(strings.Join(columns, "\t"))

	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.Split(line, "\t")) != column_count {
			// total 1831405 rows
			break
		}
		writer.WriteString(line + "\n")
	}
}
