package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dywoq/miniasm/lexer"
	"github.com/dywoq/miniasm/lexer/tokenizer"
)

func main() {	
	l, err := lexer.NewDebug(strings.NewReader("sdsd_623"), os.Stdout)
	if err != nil {
		panic(err)
	}
	
	d := tokenizer.Default{}
	d.Append(l)
	
	tokens, err := l.Do("something.miniasm")
	if err != nil {
		panic(err)
	}
	
	for _, tok := range tokens {
		fmt.Printf("tok.Literal: %v\n", tok.Literal)
	}
}
