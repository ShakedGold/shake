package parser

import (
	"fmt"
	"os"
	"shake/lexer"
	"strconv"
)

type Type int

const (
	EmptyType Type = iota
	IntType   Type = iota
	FloatType Type = iota
	CharType  Type = iota
)

type Node interface {
	Run() Node
}

type VariableDeclarationNode struct {
	Name  string
	Type  string
	Value Node
}

type FunctionDeclarationNode struct {
	Name       string
	Return     Type
	Decorator  *Decorator
	Parameters []VariableDeclarationNode
	Body       []Node
}

type Decorator struct {
	Name string
}

type IntNode struct {
	Value int
}

func (i IntNode) Run() Node {
	return i
}

type BinaryOperationNode struct {
	Operator string
	Left     Node
	Right    Node
}

func (b BinaryOperationNode) Run() Node {
	var leftNum, rightNum, result IntNode
	var ok bool
	if leftNum, ok = b.Left.(IntNode); !ok {
		fmt.Fprintln(os.Stderr, "Runtime Error: not an int on the left side on the + operation")
		return nil
	}
	if rightNum, ok = b.Right.(IntNode); !ok {
		fmt.Fprintln(os.Stderr, "Runtime Error: not an int on the right side on the + operation")
		return nil
	}

	switch b.Operator {
	case "+":
		result.Value = leftNum.Value + rightNum.Value
	case "-":
		result.Value = leftNum.Value - rightNum.Value
	case "/":
		result.Value = leftNum.Value / rightNum.Value
	case "*":
		result.Value = leftNum.Value * rightNum.Value
	}

	return result
}

type LiteralNode struct {
	Value string
}

type IdentifierNode struct {
	Name string
}

func (i IdentifierNode) Run() Node {
	return nil
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

type ReturnNode struct {
	Result Node
}

type ImportNode struct {
	Path string
}

type GlobalScopeNode interface {
	IsGlobalScopeNode()
}

func (n ImportNode) IsGlobalScopeNode()              {}
func (n FunctionDeclarationNode) IsGlobalScopeNode() {}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) nextToken() lexer.Token {
	token := p.tokens[p.pos]
	p.pos++
	return token
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) peekAhead(offset int) lexer.Token {
	if p.pos+offset < len(p.tokens) {
		return p.tokens[p.pos+offset]
	}
	return lexer.Token{Type: lexer.TokenUnknown, Value: ""}
}

func (p *Parser) expect(value string) error {
	token := p.peek()
	if token.Value != value {
		return fmt.Errorf("Expected '%s', got '%s'", value, token.Value)
	}
	return nil
}

func (p *Parser) parseStatement() (Node, error) {
	token := p.peek()
	var node Node

	for p.hasMoreTokens() {
		switch {
		case p.peekAhead(1).Type == lexer.TokenOperation:
			left, err := strconv.ParseInt(p.nextToken().Value, 10, 32)
			if err != nil {
				return nil, err
			}
			operation := p.nextToken()
			right, err := strconv.ParseInt(p.nextToken().Value, 10, 32)
			if err != nil {
				return nil, err
			}
			return BinaryOperationNode{
				Operator: operation.Value,
				Left:     IntNode{Value: int(left)},
				Right:    IntNode{Value: int(right)},
			}, nil
		case token.Type == lexer.TokenNumber:
			val, err := strconv.ParseInt(p.nextToken().Value, 10, 32)
			if err != nil {
				return nil, err
			}
			return IntNode{
				Value: int(val),
			}, nil
		case token.Type == lexer.TokenKeyword && token.Value == "if":
			// return p.parseIfStatement()
			return nil, nil
		case token.Type == lexer.TokenKeyword && token.Value == "return":

		case token.Type == lexer.TokenIdentifier:
			next := p.peekAhead(1)
			if next.Type == lexer.TokenOperation && next.Value == "=" {
				// return p.parseAssignment()
				return nil, nil
			} else {
				return IdentifierNode{
					Name: token.Value,
				}, nil
			}
		case token.Type == lexer.TokenNumber:
			p.nextToken()
			continue
		default:
			return nil, fmt.Errorf("Unexpected token in statement: %s", token.Value)
		}
	}
	return node, nil
}

func (p *Parser) parseScope() ([]Node, error) {
	numOfRemainingParenthesis := 0
	var nodes []Node
	for {
		token := p.nextToken()
		if !p.hasMoreTokens() || p.peek().Type == lexer.TokenPunctuation && p.peek().Value == "}" && numOfRemainingParenthesis == 0 {
			break
		}
		if token.Value == "{" {
			numOfRemainingParenthesis++
		} else if token.Value == "}" {
			numOfRemainingParenthesis--
		}

		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}
	p.expect("}") // Consume '}'
	return nodes, nil
}

func (p *Parser) parseVariableDeclaration() (*VariableDeclarationNode, error) {
	// Parse parameter name
	paramNameToken := p.nextToken()
	if paramNameToken.Type != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Expected parameter name, got %s", paramNameToken.Value)
	}

	// Parse parameter type
	paramTypeToken := p.nextToken()
	if paramTypeToken.Type != lexer.TokenIdentifierType {
		return nil, fmt.Errorf("Expected parameter type, got %s", paramTypeToken.Value)
	}

	return &VariableDeclarationNode{
		Name:  paramNameToken.Value,
		Type:  paramTypeToken.Value,
		Value: nil,
	}, nil
}

func (p *Parser) parseParameterList() ([]VariableDeclarationNode, error) {
	var parameters []VariableDeclarationNode

	for {
		// Check for closing parenthesis (end of parameter list)
		if p.peek().Value == ")" {
			break
		}

		varDeclaration, err := p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}

		// Add the parameter to the list
		parameters = append(parameters, *varDeclaration)

		// Check for a comma (optional) or closing parenthesis
		if p.peek().Value == "," {
			p.nextToken() // Consume the comma (,)
		} else if p.peek().Value != ")" {
			return nil, fmt.Errorf("Unexpected token in parameter list: %s", p.peek().Value)
		}
	}

	return parameters, nil
}

func (p *Parser) hasMoreTokens() bool {
	return p.pos < len(p.tokens)
}

func (p *Parser) parseFunctionDeclaration() (*FunctionDeclarationNode, error) {
	var nameToken lexer.Token
	var parameters []VariableDeclarationNode
	var err error
	var returnType Type = EmptyType

	// Optional name
	if p.peek().Type == lexer.TokenIdentifier {
		nameToken = p.nextToken() // Consume name
	}

	// Arguments
	err = p.expect("(")
	if err != nil {
		return nil, err
	}
	p.nextToken() // Consume '('
	parameters, err = p.parseParameterList()
	if err != nil {
		return nil, err
	}
	err = p.expect(")") // Consume ')'
	if err != nil {
		return nil, err
	}
	p.nextToken()

	// Optional return type
	if p.peek().Type == lexer.TokenIdentifierType {
		token := p.nextToken()
		switch token.Value {
		case "int":
			returnType = IntType
		case "float":
			returnType = FloatType
		case "char":
			returnType = CharType
		}
	}

	// Function body
	err = p.expect("{") // Consume '{'
	if err != nil {
		return nil, err
	}
	p.nextToken()
	p.nextToken()               // Consume 'NL'
	body, err := p.parseScope() // Parse function body
	if err != nil {
		return nil, err
	}
	err = p.expect("}")
	if err != nil {
		return nil, err
	}
	p.nextToken()
	if p.hasMoreTokens() {
		p.nextToken() // Consume 'NL'
	}

	return &FunctionDeclarationNode{
		Name:       nameToken.Value,
		Parameters: parameters,
		Body:       body,
		Return:     returnType,
	}, nil
}

func (p *Parser) ParseGlobalScope() ([]GlobalScopeNode, error) {
	program := []GlobalScopeNode{}

	var decorator *Decorator
	hasEntryDecorator := false

	for p.hasMoreTokens() {
		// functions
		if (p.peek().Value == "(" && p.peekAhead(3).Type != lexer.TokenIdentifier) || (p.peek().Type == lexer.TokenIdentifier && p.peekAhead(1).Value == "(") {
			// Function declaration (with or without a name)
			funcNode, err := p.parseFunctionDeclaration()
			if err != nil {
				panic(err)
			}
			if decorator != nil {
				funcNode.Decorator = decorator
				decorator = nil
			}
			program = append(program, funcNode)
			// decorators
		} else if p.peek().Value == "(" && p.peekAhead(1).Type == lexer.TokenIdentifier && p.peekAhead(2).Value == ")" {
			p.nextToken()
			token := p.nextToken()
			decorator = &Decorator{
				Name: token.Value,
			}
			if token.Value == "entry" {
				hasEntryDecorator = true
			}
			p.nextToken()
		} else {
			decorator = nil
			p.nextToken()
		}
	}
	if !hasEntryDecorator {
		return nil, fmt.Errorf("Error: no entry decorator '(entry)'")
	}
	return program, nil
}
