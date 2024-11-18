package main

import (
	"fmt"
	"shake/lexer"
)

func main() {
	source := "31 + 69"

	fmt.Println(lexer.Lex(source))
}
