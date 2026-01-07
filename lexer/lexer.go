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

package lexer

import (
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/dywoq/miniasm/lexer/tokenizer"
	"github.com/dywoq/miniasm/token"
)

type Lexer struct {
	// base
	r     io.Reader
	bytes []byte
	on    atomic.Bool

	// debug
	debugW      io.Writer
	debugOn     atomic.Bool
	debugLogger *log.Logger

	// tokenizers
	tokenizers []tokenizer.Tokenizer

	// mutex
	mu sync.Mutex

	// data
	filename string
	position *token.Position
}

// New creates a new instance of Lexer with debugging automatically turned off.
// The function tries to read bytes from r, and save it in Lexer.
//
// Returns an error if any failure is encountered.
func New(r io.Reader) (*Lexer, error) {
	l, err := newBase(r)
	if err != nil {
		return nil, err
	}
	l.debugOn.Store(false)
	l.debugW = nil
	return l, nil
}

// NewDebug creates a new instance of Lexer.
// It works the same as New, but the only difference is that NewDebug requires debug writer,
// and automatically turns debug mode on.
//
// Returns an error if any failure is encountered.
func NewDebug(r io.Reader, w io.Writer) (*Lexer, error) {
	l, err := newBase(r)
	if err != nil {
		return nil, err
	}
	l.debugOn.Store(true)
	l.debugW = w
	l.debugLogger = log.New(l.debugW, "", log.Default().Flags())
	return l, nil
}

func newBase(r io.Reader) (*Lexer, error) {
	l := &Lexer{}
	l.r = r
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l.bytes = bytes
	l.mu = sync.Mutex{}
	l.on.Store(false)
	l.position = &token.Position{}
	return l, nil
}

// SetReader sets a new reader, which updates the underlying bytes.
//
// Panics if the lexer is currently working.
//
// Returns an error if reading fails.
func (l *Lexer) SetReader(r io.Reader) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.on.Load() {
		panic("lexer is on, can't set reader")
	}
	l.r = r
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	l.bytes = bytes
	return nil
}

// DebugSetWriter sets a new debugging writer.
// Panics if the lexer is currently working.
func (l *Lexer) DebugSetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.on.Load() {
		panic("lexer is on, can't set debug writer")
	}
	l.debugW = w
}

// DebugSetMode sets a debugging mode to b.
// Panics if the lexer is currently working.
func (l *Lexer) DebugSetMode(b bool) {
	if l.on.Load() {
		panic("lexer is on, can't set debug mode")
	}
	l.debugOn.Store(b)
}

// DebugOn returns true if debugging is on.
func (l *Lexer) DebugOn() bool {
	return l.debugOn.Load()
}

// AppendTokenizer appends a new tokenizer to the lexer.
// Panics if the lexer is currently working.
func (l *Lexer) AppendTokenizer(t tokenizer.Tokenizer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.on.Load() {
		panic("lexer is on, can't append tokenizer")
	}
	l.tokenizers = append(l.tokenizers, t)
}

// implements tokenizer.Context
type context struct {
	l *Lexer
}

func (c *context) Eof() bool {
	return c.l.position.Position >= len(c.l.bytes)
}

func (c *context) Current() byte {
	if c.Eof() {
		return 0
	}
	return c.l.bytes[c.l.position.Position]
}

func (c *context) Advance() {
	if c.Eof() {
		return
	}
	c.l.position.Position++
	if cur := c.Current(); cur != 0 && cur == '\n' {
		c.l.position.Line++
		c.l.position.Column = 1
	} else {
		c.l.position.Column++
	}
}

func (c *context) Slice(start, end int) (string, error) {
	switch {
	case start > end:
		return "", c.l.makeError(fmt.Sprintf("Start %v is higher than end %v", start, end))
	case start < 0:
		return "", c.l.makeError(fmt.Sprintf("Start %v is negative", start))
	case end > len(c.l.bytes):
		return "", c.l.makeError(fmt.Sprintf("End %v is out of bounds", end))
	}
	return string(c.l.bytes[start:end]), nil
}

func (c *context) Position() *token.Position {
	return c.l.position
}

func (c *context) DebugPrintf(format string, a ...any) {
	if c.l.DebugOn() {
		c.l.debugLogger.Printf(format, a...)
	}
}

func (c *context) DebugPrint(a ...any) {
	if c.l.DebugOn() {
		c.l.debugLogger.Print(a...)
	}
}

func (c *context) DebugPrintln(a ...any) {
	if c.l.DebugOn() {
		c.l.debugLogger.Println(a...)
	}
}

// Do starts lexer and runs tokenizers, printing debug messages
// if debug mode is on. 
// 
// Does nothing if there are no set tokenizers.
//
// Returns an error if tokenizer failed to transform
// input into a token.
func (l *Lexer) Do(filename string) ([]*token.Token, error) {
	c := &context{l}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.filename = filename

	c.DebugPrintln("Starting lexer...")
	l.on.Store(true)
	defer func() {
		c.DebugPrintln("Lexer ended")
		l.on.Store(false)
	}()

	if len(l.tokenizers) == 0 {
		c.DebugPrintln("No tokenizers detected")
		return nil, nil
	}

	tokens := []*token.Token{}
	for !c.Eof() {
		tok, err := l.tokenize(c)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}

func (l *Lexer) makeError(err string) error {
	return fmt.Errorf("%v (at %v:%v:%v)", err, l.filename, l.position.Line, l.position.Column)
}

func (l *Lexer) tokenize(c *context) (*token.Token, error) {
	for _, tokenizer := range l.tokenizers {
		c.DebugPrintln("Trying to tokenize...")
		tok, noMatch, err := tokenizer(c)
		if err != nil {
			c.DebugPrintln("Encountered an error when tokenizing")
			return nil, err
		}
		if noMatch {
			c.DebugPrintln("No match, trying other tokenizer...")
			continue
		}
		c.DebugPrintln("Successfully tokenized, returning token")
		return tok, err
	}
	return nil, l.makeError("Unknown character")
}
