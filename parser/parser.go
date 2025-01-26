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

type NodeScopedStatement interface{}

type NodeScope struct {
	statements  []NodeScopedStatement
	identifiers map[string]*NodeTermIdentifier
	returnType  types.Type
}
type NodeProgram struct {
	NodeScope
	CurrentScope *NodeScope
}

type Parser struct {
	tokens  *queue.Queue[lexer.Token]
	program *NodeProgram
}

func NewParser(tokens *queue.Queue[lexer.Token]) *Parser {
	return &Parser{
		tokens: tokens,
		program: &NodeProgram{
			NodeScope: NodeScope{
				statements:  []NodeScopedStatement{},
				identifiers: make(map[string]*NodeTermIdentifier),
			},
		},
	}
}

func (p *Parser) ParseProgram() (*NodeProgram, error) {
	token, err := p.tokens.TryPop()
	for err == nil {
		if token.Type != lexer.TokenKeyword {
			return nil, ExpectedError("keywords - `fn/import`", token.LineNumber)
		}
		switch token.Value {
		case "fn":
			function, err := p.parseFunction()
			if err != nil {
				return nil, err
			}
			p.program.statements = append(p.program.statements, function)
		// TODO: imports
		case "import":
		default:
			return nil, ExpectedError("keywords - `fn/import`", token.LineNumber)
		}
		token, err = p.tokens.TryPop()
	}
	return p.program, nil
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
