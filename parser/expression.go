package parser

import "shake/types"

// TODO: Add expressions
type NodeExpression struct {
	type_ types.Type
}

func (p *Parser) parseExpression() (*NodeExpression, error) {
	p.tokens.Pop()
	return &NodeExpression{
		type_: types.TypeInt32,
	}, nil
}
