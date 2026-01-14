package mini

import (
	"github.com/dywoq/miniasm/pkg/ast"
	"github.com/dywoq/miniasm/pkg/token"
)

// Default contains default MiniASM parsers.
type Default struct {
}

func (d *Default) Append(a Appender) {
	a.AppendParser(d.TopLevel)
}

func (d *Default) TopLevel(c Context) (ast.Node, bool, error) {
	identifier, ok := c.ExpectKind(token.Identifier)
	if !ok {
		return nil, true, nil
	}
	expr, err := d.Expression(c)
	if err != nil {
		return nil, false, err
	}
	return &ast.TopLevel{Identifier: identifier.Literal, Expression: expr}, false, nil
}

func (d *Default) Expression(c Context) (ast.Node, error) {
	tok := c.Current()

	switch tok.Kind {
	case token.Number, token.Char, token.String:
		return d.Variable(c)
	default:
		return nil, c.NewError("Unknown token", tok.Position)
	}
}

func (d *Default) Variable(c Context) (ast.Node, error) {
	expect := []token.Kind{
		token.Char,
		token.Number,
		token.String,
	}
	var lastTok *token.Token
	failed := true
	for _, expected := range expect {
		tok, ok := c.ExpectKind(expected)
		if !ok {
			lastTok = tok
			continue
		} else {
			failed = false
			lastTok = tok
			break
		}
	}
	if failed && lastTok != nil {
		return nil, c.NewError("Expected Char, Number or String", lastTok.Position)
	}
	return &ast.Value{Literal: lastTok.Literal, Kind: lastTok.Kind}, nil
}
