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
// github:kevindamm/wits-osn/enum_type.go

package osn

import "fmt"

// Dependent type for EnumTable, the shape that all enums are assumed to have.
//
// Within this project, all enum-compatible values must satisfy this interface.
// An underlying type of uint8 (anything larger should be promoted to a Table)
// and an ability to be converted into a string.  The string equivalent will be
// used as the human-readable format, while elsewhere the integer value is used.
//
// It is assumed that the zero value for enums has UNKNOWN-equivalent semantics.
type EnumType interface {
	~uint8
	fmt.Stringer

	// Returns true if the represented value is valid;
	// the uint8(0) value ("UNKNOWN") should also be considered valid.
	IsValid() bool
}
