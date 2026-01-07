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

package lexer_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dywoq/miniasm/lexer"
	"github.com/dywoq/miniasm/lexer/tokenizer"
	"github.com/dywoq/miniasm/token"
)

func testTokenizer(matchByte byte, tokKind token.Kind) tokenizer.Tokenizer {
	return func(c tokenizer.Context) (*token.Token, bool, error) {
		if c.Eof() {
			return nil, false, nil
		}
		if c.Current() != matchByte {
			return nil, true, nil
		}
		pos := c.Position()
		c.Advance()
		posCopy := &token.Position{
			Line:     pos.Line,
			Column:   pos.Column,
			Position: pos.Position,
		}
		return token.New(string(matchByte), tokKind, posCopy), false, nil
	}
}

func whitespaceTokenizer() tokenizer.Tokenizer {
	return func(c tokenizer.Context) (*token.Token, bool, error) {
		if c.Eof() {
			return nil, false, nil
		}
		cur := c.Current()
		if cur != ' ' && cur != '\n' && cur != '\t' {
			return nil, true, nil
		}
		pos := c.Position()
		c.Advance()
		posCopy := &token.Position{
			Line:     pos.Line,
			Column:   pos.Column,
			Position: pos.Position,
		}
		return token.New(string(cur), token.Separator, posCopy), false, nil
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   io.Reader
		wantErr bool
	}{
		{"valid input", strings.NewReader("hello"), false},
		{"empty input", strings.NewReader(""), false},
		{"multiline input", strings.NewReader("hello\nworld"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.New(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if l == nil && !tt.wantErr {
				t.Error("New() returned nil lexer")
			}
		})
	}
}

func TestNewDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	l, err := lexer.NewDebug(strings.NewReader("test"), buf)
	if err != nil {
		t.Fatalf("NewDebug() error = %v", err)
	}
	if !l.DebugOn() {
		t.Error("NewDebug() should have debug mode on")
	}
	if l.DebugOn() && buf == nil {
		t.Error("NewDebug() should have debug writer set")
	}
}

func TestSetReader(t *testing.T) {
	l, err := lexer.New(strings.NewReader("initial"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	err = l.SetReader(strings.NewReader("updated"))
	if err != nil {
		t.Errorf("SetReader() error = %v", err)
	}
}

func TestSetReaderWhileOn(t *testing.T) {
	l, err := lexer.New(strings.NewReader("t"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('t', token.Identifier))

	done := make(chan struct{})
	go func() {
		_, _ = l.Do("test")
		close(done)
	}()

	<-done

	err = l.SetReader(strings.NewReader("new"))
	if err != nil {
		t.Errorf("SetReader() should succeed after Do() completes, got error: %v", err)
	}
}

func TestDebugSetWriter(t *testing.T) {
	l, err := lexer.New(strings.NewReader("test"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	buf := &bytes.Buffer{}
	l.DebugSetWriter(buf)
}

func TestDebugSetMode(t *testing.T) {
	l, err := lexer.New(strings.NewReader("test"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.DebugSetMode(true)
	if !l.DebugOn() {
		t.Error("DebugSetMode(true) should enable debug mode")
	}

	l.DebugSetMode(false)
	if l.DebugOn() {
		t.Error("DebugSetMode(false) should disable debug mode")
	}
}

func TestAppendTokenizer(t *testing.T) {
	l, err := lexer.New(strings.NewReader("test"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('t', token.Identifier))
	l.AppendTokenizer(testTokenizer('e', token.Identifier))
}

func TestDo(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		tokenizers     []tokenizer.Tokenizer
		expectedTokens int
		wantErr        bool
	}{
		{
			name:           "no tokenizers",
			input:          "test",
			tokenizers:     nil,
			expectedTokens: 0,
			wantErr:        false,
		},
		{
			name:           "single matching tokenizer",
			input:          "xxx",
			tokenizers:     []tokenizer.Tokenizer{testTokenizer('x', token.Identifier)},
			expectedTokens: 3,
			wantErr:        false,
		},
		{
			name:           "multiple tokenizers",
			input:          "ab",
			tokenizers:     []tokenizer.Tokenizer{testTokenizer('a', token.Identifier), testTokenizer('b', token.Identifier)},
			expectedTokens: 2,
			wantErr:        false,
		},
		{
			name:           "empty input with tokenizers",
			input:          "",
			tokenizers:     []tokenizer.Tokenizer{testTokenizer('x', token.Identifier)},
			expectedTokens: 0,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.New(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			for _, tok := range tt.tokenizers {
				l.AppendTokenizer(tok)
			}

			tokens, err := l.Do("test")
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tokens) != tt.expectedTokens {
				t.Errorf("Do() returned %d tokens, want %d", len(tokens), tt.expectedTokens)
			}
		})
	}
}

func TestDoWithWhitespace(t *testing.T) {
	l, err := lexer.New(strings.NewReader("a b c"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('a', token.Identifier))
	l.AppendTokenizer(testTokenizer('b', token.Identifier))
	l.AppendTokenizer(testTokenizer('c', token.Identifier))
	l.AppendTokenizer(whitespaceTokenizer())

	tokens, err := l.Do("test")
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(tokens) != 5 {
		t.Errorf("Do() returned %d tokens, want 5", len(tokens))
	}
}

func TestDoUnknownCharacter(t *testing.T) {
	l, err := lexer.New(strings.NewReader("#"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('a', token.Identifier))

	_, err = l.Do("test")
	if err == nil {
		t.Error("Do() should error on unknown character")
	}
}

func TestDoMultiline(t *testing.T) {
	l, err := lexer.New(strings.NewReader("a\nb\nc"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('a', token.Identifier))
	l.AppendTokenizer(testTokenizer('b', token.Identifier))
	l.AppendTokenizer(testTokenizer('c', token.Identifier))
	l.AppendTokenizer(whitespaceTokenizer())

	tokens, err := l.Do("test")
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(tokens) != 5 {
		t.Errorf("Do() returned %d tokens, want 5", len(tokens))
	}

	if tokens[2].Literal != "b" {
		t.Errorf("Token 2 should be 'b', got %q", tokens[2].Literal)
	}
}

func TestTokenPosition(t *testing.T) {
	l, err := lexer.New(strings.NewReader("ab"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('a', token.Identifier))
	l.AppendTokenizer(testTokenizer('b', token.Identifier))

	tokens, err := l.Do("test")
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if tokens[0].Position.Position != 1 {
		t.Errorf("First token position = %d, want 1", tokens[0].Position.Position)
	}

	if tokens[1].Position.Position != 2 {
		t.Errorf("Second token position = %d, want 2", tokens[1].Position.Position)
	}
}

func TestDoConcurrentSafety(t *testing.T) {
	l, err := lexer.New(strings.NewReader("abc"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	l.AppendTokenizer(testTokenizer('a', token.Identifier))

	done := make(chan struct{})
	go func() {
		_, _ = l.Do("test")
		close(done)
	}()

	<-done
}

func TestDoWithAllSeparators(t *testing.T) {
	input := ";,[]{}"
	l, err := lexer.New(strings.NewReader(input))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	for _, sep := range token.Separators {
		l.AppendTokenizer(testTokenizer(sep[0], token.Separator))
	}

	tokens, err := l.Do("test")
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(tokens) != len(input) {
		t.Errorf("Do() returned %d tokens, want %d", len(tokens), len(input))
	}
}
