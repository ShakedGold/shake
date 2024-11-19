package parser

import (
	"fmt"
	"shake/lexer"
)

type Node interface{}

type VariableDeclarationNode struct {
	Name  string
	Type  string
	Value Node
}

type FunctionDeclarationNode struct {
	Name       string
	Parameters []VariableDeclarationNode
	Body       []Node
}

type BinaryOperationNode struct {
	Operator string
	Left     Node
	Right    Node
}

type LiteralNode struct {
	Value string
}

type IdentifierNode struct {
	Name string
}

type AssignmentNode struct {
	Name  string
	Value Node
}

type IfNode struct {
	Condition Node
	Body      []Node
	ElseBody  []Node
}

type ImportNode struct {
	Path string
}

type GlobalScopeNode interface {
	IsGlobalScopeNode()
}

func (n *ImportNode) IsGlobalScopeNode()              {}
func (n *FunctionDeclarationNode) IsGlobalScopeNode() {}

type ProgramNode struct {
	GlobalStatements []GlobalScopeNode
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func (p *Parser) nextToken() lexer.Token {
	token := p.tokens[p.pos]
	p.pos++
	return token
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) hasMoreTokens() bool {
	return p.pos < len(p.tokens)
}

func (p *Parser) ParseProgram() *ProgramNode {
	program := &ProgramNode{}

	for p.hasMoreTokens() {
		token := p.peek()
		switch {
		case token.Type == lexer.TokenImport:
			// program.GlobalStatements = append(program.GlobalStatements, p.ParseImport())
		case token.Type == lexer.TokenFunction:
			program.GlobalStatements = append(program.GlobalStatements, p.ParseFunction())
		default:
			panic(fmt.Sprintf("Invalid statement in global scope: %s", token.Value))
		}
	}

	return program
}
