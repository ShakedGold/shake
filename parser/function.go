package parser

import (
	"encoding/json"
	"fmt"
	"shake/lexer"
)

type NodeFunction struct {
	scope *NodeScope
	name  string
}

func (nf NodeFunction) MarshalJSON() ([]byte, error) {
	return json.Marshal(nf.name)
}

func (p *Parser) parseFunction() (*NodeFunction, error) {
	// expected identifier: `main/add`
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenIdentifier})
	if err != nil {
		return nil, err
	}
	funcIdentifier := p.tokens.Pop()
	nodeFunction := &NodeFunction{
		name: funcIdentifier.Value,
	}

	// expected `(`
	token, err = p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenPunctuation, Value: "("})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	// expected `)`
	token, err = p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenPunctuation, Value: ")"})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	// expected return type
	token, err = p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenIdentifierType})
	if err != nil {
		return nil, ExpectedError(fmt.Sprintf("Type got: %s", token.Value), token.LineNumber)
	}

	p.tokens.Pop()

	// parse scope``
	scope, err := p.parseScope()
	if err != nil {
		return nil, err
	}
	nodeFunction.scope = scope

	return nodeFunction, nil
}
