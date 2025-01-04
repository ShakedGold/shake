package parser

import "shake/lexer"

func (p *Parser) parseScope() (*NodeScope, error) {
	// expect `{`
	err := expectToken(p.tokens.Peek(0), lexer.Token{Type: lexer.TokenPunctuation, Value: "{"})
	if err != nil {
		return nil, err
	}
	p.tokens.Pop()

	scope := &NodeScope{
		statements: []NodeScopedStatement{},
	}

	// parse statements until {
	for {
		currTokenOpt := p.tokens.Peek(0)
		if !currTokenOpt.Exists() {
			return nil, ExpectedError("`{` but found nothing", 0)
		}
		currToken := currTokenOpt.Value()
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
	return scope, nil
}
