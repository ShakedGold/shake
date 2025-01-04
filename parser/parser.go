package parser

import (
	"fmt"
	"shake/lexer"
	"shake/optional"
	"shake/queue"
	"shake/types"
)

type NodeTerm interface{}
type NodeGlobalStatements interface{}
type NodeScopedStatement interface{}
type NodeExpression interface{}

type NodeAssignment struct {
	name       string
	type_      types.Type
	expression NodeExpression
}
type NodeIdentifier struct {
	token lexer.Token
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

	for p.tokens.Peek(0).Exists() {
		token := p.tokens.Pop()
		if token.Type != lexer.TokenKeyword {
			return nil, ExpectedError("keywords - (fn/import)", token.LineNumber)
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
			return nil, ExpectedError("keywords - (fn/import)", token.LineNumber)
		}
	}
	return program, nil
}

func Error(reason string, line uint64) error {
	return fmt.Errorf("[Parser Error]: %s, At line: %d", reason, line)
}
func ExpectedError(reason string, line uint64) error {
	return Error("Expected "+reason, line)
}
func expectToken(optToken optional.Optional[lexer.Token], token lexer.Token) error {
	if !optToken.Exists() {
		return ExpectedError(fmt.Sprintf("`%s` but didn't find anything", token.Value), 0)
	}
	currToken := optToken.Value()

	if token.Value == "" && token.Type != currToken.Type {
		return ExpectedError(fmt.Sprintf("%s but found: %s", token.Type, currToken.Value), currToken.LineNumber)
	} else if token.Value != "" && (token.Type != currToken.Type || token.Value != currToken.Value) {
		return ExpectedError(fmt.Sprintf("%s but found: %s", token.Value, currToken.Value), currToken.LineNumber)
	}

	return nil
}
