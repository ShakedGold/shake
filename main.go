package main

import (
	"bytes"
	"fmt"
	"os"
	"shake/lexer"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Not enough arguments")
		os.Exit(1)
	}

	programSourcePath := args[1]
	programSource, err := os.ReadFile(programSourcePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read file")
		os.Exit(2)
	}
	programReader := bytes.NewReader(programSource)
	fmt.Println(lexer.Lex(programReader))
}
