package main

import (
	"fmt"
	"os"

	"github.com/dywoq/miniasm/pkg/lexer"
	"github.com/dywoq/miniasm/pkg/lexer/tokenizer"
	"github.com/dywoq/miniasm/pkg/parser"
)

func main() {
	f, err := os.Open("main.miniasm")
	if err != nil {
		panic(err)
	}

	l, err := lexer.NewDebug(f, os.Stdout)
	if err != nil {
		panic(err)
	}

	d := tokenizer.Default{}
	d.Append(l)

	tokens, err := l.Do(f.Name())
	if err != nil {
		panic(err)
	}

	p := parser.NewDebug(tokens, os.Stdout)

	tree, err := p.Do(f.Name())
	if err != nil {
		panic(err)
	}

	for _, topLevel := range tree.TopLevel {
		fmt.Printf("%v", topLevel)
	}
}
