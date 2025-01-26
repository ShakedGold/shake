package parser

import (
	"fmt"
	"shake/lexer"
	"shake/types"
)

// TODO: Add expressions

/*
	{
		type: Int32
		value: {
			type: Int32
			value: 3
		}
	}
*/
type NodeExpression interface {
	GetType() types.Type
}
type NodeTerm interface {
	GetType() types.Type
}
type NodeTermInt32 struct {
	Value string
}

func (nti NodeTermInt32) GetType() types.Type {
	return types.TypeInt32
}

type NodeTermIdentifier struct {
	Type       types.Type
	Identifier string
}

func (nti NodeTermIdentifier) GetType() types.Type {
	return nti.Type
}

type NodeExpressionBinary struct {
	Left      NodeExpression
	Right     NodeExpression
	Operation string
}

func (neb NodeExpressionBinary) GetType() types.Type {
	return types.TypeInt32
}

type NodeExpressionLiteral struct {
	Type  types.Type
	Value NodeTerm
}

func (nel NodeExpressionLiteral) GetType() types.Type {
	return nel.Type
}

type NodeExpressionIdentifier struct {
	Type       types.Type
	Identifier NodeTerm
}

func (nei NodeExpressionIdentifier) GetType() types.Type {
	return nei.Type
}

func (p *Parser) parseTerm() (NodeTerm, error) {
	// current token is the term
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("token but found nothing", 0)
	}
	p.tokens.Pop()
	switch token.Type {
	case lexer.TokenIdentifier:
		// check if the identifier exists
		identifier, ok := p.program.identifiers[token.Value]
		if !ok {
			return nil, Error(fmt.Sprintf("Identifier: %s of type: %s does not exist in the current scope", token.Value, token.Type), token.LineNumber)
		}
		return identifier, nil
	case lexer.TokenNumber:
		return NodeTermInt32{
			Value: token.Value,
		}, nil
	default:
		return nil, ExpectedError(fmt.Sprintf("number or identifier but found: %s", token.Type), token.LineNumber)
	}
}

func (p *Parser) parseExpression() (NodeExpression, error) {
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("token but found nothing", 0)
	}

	// check if the next token is */+-
	nextToken, err := p.tokens.Peek(1)
	if err != nil {
		return nil, ExpectedError("token but found nothing", 0)
	}
	if nextToken.Type == lexer.TokenOperation {
		left, err := p.parseTerm()
		operation := p.tokens.Pop()
		if err != nil {
			return nil, ExpectedError("term but found nothing", token.LineNumber)
		}
		right, err := p.parseExpression()
		if err != nil {
			return nil, ExpectedError("expression but found nothing", token.LineNumber)
		}
		return &NodeExpressionBinary{
			Left: &NodeExpressionLiteral{
				Type:  left.GetType(),
				Value: left,
			},
			Right:     right,
			Operation: operation.Value,
		}, nil
	}

	// not an operation so parseTerm
	term, err := p.parseTerm()
	if err != nil {
		return nil, ExpectedError("term but found nothing", token.LineNumber)
	}
	switch token.Type {
	case lexer.TokenIdentifier:
		return &NodeExpressionIdentifier{
			Type:       term.GetType(),
			Identifier: term,
		}, nil
	case lexer.TokenNumber:
		return NodeExpressionLiteral{
			Type:  term.GetType(),
			Value: term,
		}, nil
	default:
		return nil, ExpectedError(fmt.Sprintf("number or identifier but found: %s", token.Type), token.LineNumber)
	}
}
