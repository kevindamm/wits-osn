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
// github:kevindamm/wits-osn/boolish.go

package osn

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Represents a boolean-like value that may be a bool, an int or a string.
// Implements JSON-codec and SQL-codec operations for natural encoding/decoding.
//
// If non-boolean, a 0 or "0" or "" value is false and any other value is true.
// When serializing, it uses the native boolean type.
type Boolish bool

func (b *Boolish) UnmarshalJSON(encoded []byte) error {
	var boolVal bool
	if err := json.Unmarshal(encoded, &boolVal); err != nil {
		*b = Boolish(boolVal)
	} else {
		var intVal int
		if err := json.Unmarshal(encoded, &intVal); err != nil {
			*b = Boolish(intVal != 0)
		} else {
			var strVal string
			if err := json.Unmarshal(encoded, &strVal); err != nil {
				*b = Boolish(!(strVal == "" || strVal == "0"))
			} else {
				return fmt.Errorf("failed to convert [%s] to Boolish type", encoded)
			}
		}
	}
	return nil
}

func (b Boolish) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

// Similar to the conversion rules above but any non-{0|1} integer is an error,
// this and string conversion are consistent with [driver.Bool] translation.
func (b *Boolish) Scan(value interface{}) error {
	if value == nil {
		*b = Boolish(false)
		return nil
	}
	if boolValue, err := driver.Bool.ConvertValue(value); err == nil {
		if v, ok := boolValue.(bool); ok {
			*b = Boolish(v)
			return nil
		}
	}
	return fmt.Errorf("failed to scan Boolish value %v", value)
}

func (b Boolish) Value() (driver.Value, error) {
	return bool(b), nil
}
