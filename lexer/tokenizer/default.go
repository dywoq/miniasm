// Copyright 2026 dywoq
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tokenizer

import (
	"slices"
	"unicode"

	"github.com/dywoq/miniasm/token"
)

// Default contains default MiniAsm tokenizers.
type Default struct {
}

// Append appends all tokenizers into a.
func (d *Default) Append(a Appender) {
	a.AppendTokenizer(d.Identifier)
	a.AppendTokenizer(d.Number)
	a.AppendTokenizer(d.Separator)
}

func (d *Default) Identifier(c Context) (*token.Token, bool, error) {
	c.DebugPrintln("Identifier(): Met a possible identifier")

	cur := c.Current()
	if cur == 0 || !unicode.IsLetter(rune(cur)) && cur != '_' {
		c.DebugPrintln("Identifier(): No match")
		return nil, true, nil
	}

	start := c.Position().Position
	for {
		cur := c.Current()
		if cur == 0 || (!unicode.IsLetter(rune(cur)) && cur != '_' && !unicode.IsDigit(rune(cur))) {
			break
		}
		c.Advance()
	}

	end := c.Position().Position
	str, err := c.Slice(start, end)
	if err != nil {
		return nil, false, err
	}
	c.DebugPrintf("Identifier(): Got %v\n", str)

	if !token.IsIdentifier(str) {
		c.DebugPrintf("Identifier(): %v is not identifier\n", str)
		for range end - start {
			c.Backward()
		}
		return nil, true, err
	}

	c.DebugPrintf("Identifier(): %v is identifier\n", str)
	return token.New(str, token.Identifier, c.Position()), false, nil
}

func (d *Default) Number(c Context) (*token.Token, bool, error) {
	c.DebugPrintln("Number(): Met a possible number")

	cur := c.Current()
	if cur == 0 || !unicode.IsDigit(rune(cur)) {
		c.DebugPrintln("Number(): No match")
		return nil, true, nil
	}

	start := c.Position().Position
	for {
		cur := c.Current()
		if cur == 0 || !unicode.IsDigit(rune(cur)) {
			break
		}
		c.Advance()
	}
	end := c.Position().Position

	str, err := c.Slice(start, end)
	if err != nil {
		return nil, false, err
	}
	c.DebugPrintf("Number(): %v is number\n", str)
	return token.New(str, token.Number, c.Position()), false, nil
}

func (d *Default) Separator(c Context) (*token.Token, bool, error) {
	c.DebugPrintln("Separator(): Met a possible separator")
	cur := c.Current()
	if cur == 0 || !slices.Contains(token.Separators, string(cur)) {
		c.DebugPrintln("Separator(): No match")
		return nil, true, nil
	}
	c.DebugPrintf("Separator(): %v is separator\n", cur)
	c.Advance()
	return token.New(string(cur), token.Separator, c.Position()), false, nil
}
