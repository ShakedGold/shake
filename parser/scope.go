package parser

import "shake/lexer"

func (p *Parser) parseScope() (*NodeScope, error) {
	// expect `{`
	token, err := p.tokens.Peek(0)
	if err != nil {
		return nil, ExpectedError("`{` but found nothing", 0)
	}
	err = expectToken(token, lexer.Token{Type: lexer.TokenPunctuation, Value: "{"})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	scope := &NodeScope{
		statements: []NodeScopedStatement{},
	}
	// set current scope
	lastScope := p.program.CurrentScope
	p.program.CurrentScope = scope

	// parse statements until {
	for {
		currToken, err := p.tokens.Peek(0)
		if err != nil {
			return nil, ExpectedError("`{` but found nothing", 0)
		}
		if currToken.Type == lexer.TokenPunctuation && currToken.Value == "}" {
			p.tokens.Pop()
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		scope.statements = append(scope.statements, statement)
	}

	// unset current scope
	p.program.CurrentScope = lastScope
	return scope, nil
}
