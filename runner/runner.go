package runner

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

type Engine string

const (
	INTERPRETER Engine = "eval"
	VM          Engine = "vm"
)

func ParseEngine(raw string) (Engine, error) {
	switch raw {
	case "vm":
		return VM, nil
	case "eval":
		return INTERPRETER, nil
	default:
		return "", fmt.Errorf("unknown engine: %s", raw)
	}
}

func ParseCode(input string) (*ast.Program, *object.Error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return nil, object.NewParserErrors(nil, p.Errors())
	}
	return program, nil
}
