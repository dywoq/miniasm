package mini

import (
	"github.com/dywoq/miniasm/pkg/ast"
	"github.com/dywoq/miniasm/pkg/debug"
	"github.com/dywoq/miniasm/pkg/token"
)

type Context interface {
	Current() *token.Token

	Advance()

	Position() int

	NewError(str string) error

	ExpectStr(str string) (*token.Token, bool)

	ExpectKind(kind token.Kind) (*token.Token, bool)

	debug.Context
}

// Parser parses tokens into AST tree node.
// Each mini parser has only one responsibility for parsing something.
// 
// Returns an error if it got what didn't expect.
type Parser func(c Context) (ast.Node, error)
