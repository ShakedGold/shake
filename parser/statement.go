package parser

import (
	"fmt"
	"shake/lexer"
	"shake/types"
)

func (p *Parser) parseStatement() (NodeScopedStatement, error) {
	// allowed statements: assignment

	// assignment
	identifierOpt := p.tokens.Peek(0)
	if !identifierOpt.Exists() {
		return nil, ExpectedError("statement but found nothing", 0)
	}
	identifier := p.tokens.Pop()

	// for now only assignments are allowed so if not an identifier we error
	if identifier.Type != lexer.TokenIdentifier {
		return nil, ExpectedError(fmt.Sprintf("identifier but found: %s", identifier.Value), identifier.LineNumber)
	}

	identifierTypeOpt := p.tokens.Peek(0)
	// TODO: type is not necessary if the expression has a type
	if !identifierTypeOpt.Exists() {
		return nil, ExpectedError("type but found nothing", 0)
	}
	typeString := p.tokens.Pop()
	identifierType := types.GetType(typeString.Value)
	if identifierType == types.TypeUnknown {
		return nil, ExpectedError(fmt.Sprintf("Type got: %s", typeString.Value), typeString.LineNumber)
	}

	// consume the `=`
	p.tokens.Pop()

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// consume the `;`
	err = expectToken(p.tokens.Peek(0), lexer.Token{Type: lexer.TokenSemicolon})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	scopedStatement := &NodeAssignment{
		name:       identifier.Value,
		type_:      identifierType,
		expression: expression,
	}

	return scopedStatement, nil
}
