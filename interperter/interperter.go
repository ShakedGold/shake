package interperter

import (
	"fmt"
	"os"
	"shake/parser"
)

type Program struct {
	globalStatements []parser.GlobalScopeNode
}

func NewProgram(globalStatements []parser.GlobalScopeNode) *Program {
	return &Program{
		globalStatements: globalStatements,
	}
}

func runMain(function *parser.FunctionDeclarationNode) int {
	for _, line := range function.Body {
		if resultNode, ok := line.Run().(parser.IntNode); ok {
			return resultNode.Value
		}
	}
	return 0
}

func (p *Program) Run() int {
	for _, statement := range p.globalStatements {
		funcDec, ok := statement.(*parser.FunctionDeclarationNode)
		if !ok {
			continue
		}
		if funcDec.Decorator == nil || funcDec.Decorator.Name != "entry" {
			continue
		}
		return runMain(funcDec)
	}
	fmt.Fprintln(os.Stderr, "Runtime Error: no main function decorator")
	return 1
}
