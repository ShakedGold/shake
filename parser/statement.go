package parser

import (
	"fmt"
	"shake/lexer"
	"shake/types"
)

type NodeAssignment struct {
	Identifier string
	Type       types.Type
	Expression *NodeExpression
}
type NodeReturn struct {
	value *NodeExpression
}

func (p *Parser) parseReturn() (*NodeReturn, error) {
	// consume the `return` keyword
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("statement but found nothing", 0)
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenKeyword, Value: "return"})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()
	// get the return value
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	currentScope := p.program.CurrentScope
	if currentScope.returnType != expression.GetType() {
		return nil, Error(fmt.Sprintf("Type of scope: %s is different from return type: %s", currentScope.returnType, expression.GetType().String()), token.LineNumber)
	}

	// consume `;`
	token, err = p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenSemicolon})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	return &NodeReturn{
		value: &expression,
	}, nil
}

func (p *Parser) parseAssignment() (*NodeAssignment, error) {
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("statement but found nothing", 0)
	}
	identifier := p.tokens.Pop()

	// for now only assignments are allowed so if not an identifier we error
	if identifier.Type != lexer.TokenIdentifier {
		return nil, ExpectedError(fmt.Sprintf("identifier but found: `%s`", identifier.Value), identifier.LineNumber)
	}

	// only consume type if exists and if not get the expression type
	identifierToken, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("type but found nothing", 0)
	}

	identifierType := types.GetType(identifierToken.Value)
	var typeString *lexer.Token
	if identifierType != types.TypeUnknown {
		typeString = p.tokens.Pop()
	}

	// consume the `=`
	p.tokens.Pop()

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if identifierType == types.TypeUnknown {
		identifierType = expression.GetType()
	}

	// check variable type and expression type match
	if identifierType != expression.GetType() {
		return nil, Error(fmt.Sprintf("Mismatched type when assigning variable %s of type %s and expression of type %s", identifier.Value, identifierType.String(), expression.GetType().String()), typeString.LineNumber)
	}

	// consume the `;`
	token, err = p.tokens.Peek(0)
	if err != nil {
		return nil, err
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenSemicolon})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	assignment := &NodeAssignment{
		Identifier: identifier.Value,
		Type:       identifierType,
		Expression: &expression,
	}

	// create the variable in the identifiers map
	p.program.identifiers[identifier.Value] = &NodeTermIdentifier{
		Type:       identifierType,
		Identifier: identifier.Value,
	}
	return assignment, nil
}

func (p *Parser) parseStatement() (NodeScopedStatement, error) {
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("token but found nothing", 0)
	}
	if token.Type == lexer.TokenKeyword && token.Value == "return" {
		return p.parseReturn()
	}
	return p.parseAssignment()
}
