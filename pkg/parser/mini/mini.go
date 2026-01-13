package mini

import (
	"github.com/dywoq/miniasm/pkg/ast"
	"github.com/dywoq/miniasm/pkg/debug"
	"github.com/dywoq/miniasm/pkg/token"
)

type Context interface {
	// IsEnd reports whether the parser has reached the end.
	IsEnd() bool

	// Current returns the current token.
	// If the parser reached the end, the function returns nil.
	Current() *token.Token

	// Advance advances to the next token.
	// If the parser reached the end, the function returns.
	Advance()

	// Position returns the current position.
	Position() int

	// NewError makes a new error with automatically inserted token position
	// and parser position.
	NewError(str string, pos *token.Position) error

	// ExpectLiteral expects a literal from the current token.
	// If the token literal doesn't satisfy lit, the function returns false,
	// otherwise, it returns true.
	//
	// The function automatically advances when returns true.
	ExpectLiteral(lit string) (*token.Token, bool)

	// ExpectKind expects a kind from the current token.
	// If the token kind doesn't satisfy kind, the function returns false,
	// otherwise, it returns true.
	//
	// The function automatically advances when returns true.
	ExpectKind(kind token.Kind) (*token.Token, bool)

	debug.Context
}

// Parser parses tokens into AST tree node.
// Each mini parser has only one responsibility for parsing something.
//
// Returns true if token doesn't match mini parser requirements.
type Parser func(c Context) (ast.Node, bool, error)

// Appended defines an interface for appending mini parsers.
type Appender interface {
	AppendParser(p Parser)
}
