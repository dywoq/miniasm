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

package parser

import (
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/dywoq/miniasm/pkg/ast"
	"github.com/dywoq/miniasm/pkg/parser/mini"
	"github.com/dywoq/miniasm/pkg/token"
)

type Parser struct {
	// base
	tokens []*token.Token
	pos    int
	on     atomic.Bool

	// debug
	debugW      io.Writer
	debugOn     atomic.Bool
	debugLogger *log.Logger

	// mini parsers
	minis []mini.Parser

	// mutex
	mu sync.Mutex

	// data
	filename string
}

func New(tokens []*token.Token) *Parser {
	p := newBase(tokens)
	return p
}

func NewDebug(tokens []*token.Token, w io.Writer) *Parser {
	p := newBase(tokens)
	p.debugW = w
	p.debugOn.Store(true)
	return p
}

func newBase(tokens []*token.Token) *Parser {
	p := &Parser{}
	p.on.Store(true)
	p.pos = 0
	p.tokens = tokens
	p.debugOn.Store(false)
	p.mu = sync.Mutex{}
	p.debugLogger = log.New(p.debugW, "", log.Default().Flags())
	p.filename = ""
	return nil
}

func (p *Parser) SetTokens(tokens []*token.Token) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.on.Load() {
		panic("parser is on, can't set tokens")
	}
	p.tokens = tokens
}

// DebugSetWriter sets a new debugging writer.
// Panics if the parser is currently working.
func (p *Parser) DebugSetWriter(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.on.Load() {
		panic("parser is on, can't set debug writer")
	}
	p.debugW = w
}

// DebugSetMode sets a debugging mode to b.
// Panics if the parser is currently working.
func (p *Parser) DebugSetMode(b bool) {
	if p.on.Load() {
		panic("parser is on, can't set debug mode")
	}
	p.debugOn.Store(b)
}

// DebugOn returns true if debugging is on.
func (p *Parser) DebugOn() bool {
	return p.debugOn.Load()
}

// implements mini.Context
type context struct {
	p *Parser
}

func (c *context) Current() *token.Token {
	if c.p.pos >= len(c.p.tokens) {
		return nil
	}
	return c.p.tokens[c.p.pos]
}

func (c *context) Advance() {
	if c.p.pos >= len(c.p.tokens) {
		return
	}
	c.p.pos++
}

func (c *context) Position() int {
	return c.p.pos
}

func (c *context) NewError(str string, pos *token.Position) error {
	return c.p.makeError(str, pos)
}

func (c *context) ExpectLiteral(lit string) (*token.Token, bool) {
	tok := c.p.tokens[c.p.pos]
	if tok.Literal != lit {
		return nil, false
	}
	c.Advance()
	return tok, true
}

func (c *context) ExpectKind(kind token.Kind) (*token.Token, bool) {
	tok := c.p.tokens[c.p.pos]
	if tok.Kind != kind {
		return nil, false
	}
	c.Advance()
	return tok, true
}

func (c *context) DebugPrintf(format string, a ...any) {
	if c.p.DebugOn() {
		c.p.debugLogger.Printf(format, a...)
	}
}

func (c *context) DebugPrint(a ...any) {
	if c.p.DebugOn() {
		c.p.debugLogger.Print(a...)
	}
}

func (c *context) DebugPrintln(a ...any) {
	if c.p.DebugOn() {
		c.p.debugLogger.Println(a...)
	}
}

func (p *Parser) makeError(str string, pos *token.Position) error {
	return fmt.Errorf("%v (%v:%v:%v)", str, p.filename, pos.Line, pos.Column)
}

func (p *Parser) Do(filename string) (*ast.Tree, error) {
	c := &context{p}
	
	p.mu.Lock()
	defer p.mu.Unlock()

	p.filename = filename
	
	c.DebugPrintln("Starting parser...")
	p.on.Store(true)
	defer func() {
		p.on.Store(false)
		c.DebugPrintln("Parser ended")
	}()
	
	if len(p.minis) == 0 {
		c.DebugPrintln("No mini parsers detected")
		return nil, nil
	}

	topLevel := []ast.Node{}
	return &ast.Tree{
		TopLevel: topLevel,
	}, nil
}
