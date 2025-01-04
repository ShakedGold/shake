package parser

import (
	"fmt"
	"shake/lexer"
	"shake/types"
)

type NodeAssignment struct {
	name       string
	type_      types.Type
	expression *NodeExpression
}
type NodeReturn struct {
	value *NodeExpression
}

func (p *Parser) parseReturn() (*NodeReturn, error) {
	// consume the `return` keyword
	err := expectToken(p.tokens.Peek(0), lexer.Token{Type: lexer.TokenKeyword, Value: "return"})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()
	// get the return value
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// consume `;`
	err = expectToken(p.tokens.Peek(0), lexer.Token{Type: lexer.TokenSemicolon})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	return &NodeReturn{
		value: expression,
	}, nil
}

func (p *Parser) parseAssignment() (*NodeAssignment, error) {
	identifierOpt := p.tokens.Peek(0)
	if !identifierOpt.Exists() {
		return nil, ExpectedError("statement but found nothing", 0)
	}
	identifier := p.tokens.Pop()

	// for now only assignments are allowed so if not an identifier we error
	if identifier.Type != lexer.TokenIdentifier {
		return nil, ExpectedError(fmt.Sprintf("identifier but found: `%s`", identifier.Value), identifier.LineNumber)
	}

	// only consume type if exists and if not get the expression type
	identifierTypeOpt := p.tokens.Peek(0)
	if !identifierTypeOpt.Exists() {
		return nil, ExpectedError("type but found nothing", 0)
	}

	var typeString lexer.Token
	var identifierType types.Type = types.TypeUnknown
	if types.GetType(identifierTypeOpt.Value().Value) != types.TypeUnknown {
		typeString = p.tokens.Pop()
		identifierType = types.GetType(typeString.Value)
	}

	// consume the `=`
	p.tokens.Pop()

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if identifierType == types.TypeUnknown {
		identifierType = expression.type_
	}

	// check variable type and expression type match
	if identifierType != expression.type_ {
		return nil, Error(fmt.Sprintf("Mismatched type when assigning variable %s of type %s and expression of type %s", identifier.Value, identifierType.String(), expression.type_.String()), typeString.LineNumber)
	}

	// consume the `;`
	err = expectToken(p.tokens.Peek(0), lexer.Token{Type: lexer.TokenSemicolon})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	assignment := &NodeAssignment{
		name:       identifier.Value,
		type_:      identifierType,
		expression: expression,
	}
	return assignment, nil
}

func (p *Parser) parseStatement() (NodeScopedStatement, error) {
	if p.tokens.Peek(0).Exists() && p.tokens.Peek(0).Value().Type == lexer.TokenKeyword && p.tokens.Peek(0).Value().Value == "return" {
		return p.parseReturn()
	}
	return p.parseAssignment()
}
