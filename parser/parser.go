package parser

import (
	"fmt"
	"runtime/debug"
	"shake/lexer"
	"shake/options"
	"shake/queue"
	"shake/types"

	"github.com/fatih/color"
)

type NodeGlobalStatements interface{}
type NodeScopedStatement interface{}

type NodeIdentifier struct {
	token lexer.Token
	type_ types.Type
}
type NodeScope struct {
	statements []NodeScopedStatement
}
type NodeProgram struct {
	statements []NodeGlobalStatements
}
type Parser struct {
	tokens *queue.Queue[lexer.Token]
}

func NewParser(tokens *queue.Queue[lexer.Token]) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) ParseProgram() (*NodeProgram, error) {
	program := &NodeProgram{
		statements: []NodeGlobalStatements{},
	}

	var err error
	for err != nil {
		token := p.tokens.Pop()
		if token.Type != lexer.TokenKeyword {
			return nil, ExpectedError("keywords - `fn/import`", token.LineNumber)
		}
		switch token.Value {
		case "fn":
			function, err := p.parseFunction()
			if err != nil {
				return nil, err
			}
			program.statements = append(program.statements, function)
		// TODO: imports
		case "import":
		default:
			return nil, ExpectedError("keywords - `fn/import`", token.LineNumber)
		}
	}
	return program, nil
}

func Error(reason string, line uint64) error {
	if len(options.Options.Verbose) > 0 && options.Options.Verbose[0] {
		debug.PrintStack()
	}
	c := color.New(color.FgRed).Add(color.Underline)
	return fmt.Errorf("%s: %s at line: %d", c.Sprint("[Parser Error]"), reason, line)
}
func ExpectedError(reason string, line uint64) error {
	return Error("Expected "+reason, line)
}
func expectToken(currToken *lexer.Token, token lexer.Token) error {
	if token.Value == "" && token.Type != currToken.Type {
		return ExpectedError(fmt.Sprintf("%s but found: %s", token.Type, currToken.Value), currToken.LineNumber)
	} else if token.Value != "" && (token.Type != currToken.Type || token.Value != currToken.Value) {
		return ExpectedError(fmt.Sprintf("%s but found: %s", token.Value, currToken.Value), currToken.LineNumber)
	}

	return nil
}
