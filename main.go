package main

import (
	"bytes"
	"fmt"
	"os"
	"shake/lexer"
	"shake/options"
	"shake/parser"

	"github.com/jessevdk/go-flags"
)

func main() {
	flags.Parse(&options.Options)
	programSource, err := os.ReadFile(options.Options.Input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read file")
		os.Exit(1)
	}
	programReader := bytes.NewReader(programSource)
	tokens, err := lexer.Lex(programReader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not lex file")
		os.Exit(2)
	}

	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*program)
}
