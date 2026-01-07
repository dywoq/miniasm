package main

import (
	"fmt"
	"os"

	"github.com/dywoq/miniasm/lexer"
	"github.com/dywoq/miniasm/lexer/tokenizer"
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

	for _, tok := range tokens {
		fmt.Printf("%v %v\n", tok.Literal, tok.Kind)
	}
}
