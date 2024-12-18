package main

import (
	"bytes"
	"fmt"
	"os"
	"shake/interperter"
	"shake/lexer"
	"shake/parser"

	"github.com/jessevdk/go-flags"
)

var options struct {
	Verbose     []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Input       string `short:"i" long:"input" description:"input shk file"`
	Interperted bool   `short:"p" long:"interperted" description:"compile and run in interperted mode"`
}

func main() {
	flags.Parse(&options)
	programSource, err := os.ReadFile(options.Input)
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

	parse := parser.NewParser(tokens)
	globalScope, err := parse.ParseGlobalScope()
	if globalScope == nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if options.Interperted {
		program := interperter.NewProgram(globalScope)
		returnValue := program.Run()
		os.Exit(returnValue)
	}
}
