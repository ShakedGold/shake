package parser

import (
	"fmt"
	"shake/lexer"
	"shake/types"
)

type NodeFunction struct {
	scope      *NodeScope
	name       string
	returnType types.Type
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
	if err == nil {
		returnType := p.tokens.Pop()
		nodeFunction.returnType = types.GetType(returnType.Value)
		if nodeFunction.returnType == types.TypeUnknown {
			return nil, ExpectedError(fmt.Sprintf("Type got: %s", returnType.Value), returnType.LineNumber)
		}
	}

	// parse scope``
	scope, err := p.parseScope()
	if err != nil {
		return nil, err
	}
	nodeFunction.scope = scope

	return nodeFunction, nil
}
