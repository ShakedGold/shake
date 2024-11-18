package lexer

import (
	"regexp"
	"shake/queue"
)

type TokenType int

const (
	TokenNone TokenType = iota
	TokenNumber
	TokenOperation
	TokenSymbol
)

var tokenNames = map[TokenType]string{
	TokenNone:      "None",
	TokenNumber:    "Number",
	TokenOperation: "Operation",
	TokenSymbol:    "Symbol",
}

func (tt TokenType) String() string {
	return tokenNames[tt]
}

type Token struct {
	Type  TokenType
	Value string
}

func Scan(queue *queue.Queue[string], first string, allowed string) string {
	number := first
	for queue.Size() != 0 {
		if next := queue.Pop(); regexp.MustCompile(allowed).MatchString(next) {
			number += next
		} else {
			break
		}
	}

	return number
}

func Lex(source string) []Token {
	sourceQueue := queue.NewQueue[string]()
	target := []Token{}

	// push all to the queue
	for _, char := range source {
		sourceQueue.Push(string(char))
	}

	for sourceQueue.Size() != 0 {
		char := sourceQueue.Pop()

		if regexp.MustCompile("[.0-9]").MatchString(char) {
			number := Scan(sourceQueue, char, "[.0-9]")
			token := Token{
				Type:  TokenNumber,
				Value: number,
			}
			target = append(target, token)
		} else if regexp.MustCompile(`[\+\-\*\/\=]`).MatchString(char) {
			token := Token{
				Type:  TokenOperation,
				Value: char,
			}
			target = append(target, token)
		} else if regexp.MustCompile(`[_a-zA-Z]`).MatchString(char) {
			symbol := Scan(sourceQueue, char, "[_a-zA-Z]")
			token := Token{
				Type:  TokenSymbol,
				Value: symbol,
			}
			target = append(target, token)
		}
	}

	return target
}
