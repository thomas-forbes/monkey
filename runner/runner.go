package runner

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"time"
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

func RunProgram(engine Engine, input string, session Session) (object.Object, Session, time.Duration) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return object.NewParserErrors(nil, p.Errors()), session, 0
	}

	var result object.Object
	var duration time.Duration

	switch engine {
	case INTERPRETER:
		result, session, duration = runEval(program, session)
	case VM:
		result, session, duration = runVM(program, session)
	default:
		return &object.Error{Message: fmt.Sprintf("unknown engine: %s", engine)}, session, 0
	}

	if result == nil {
		return object.NULL, session, duration
	}

	return result, session, duration
}

func runEval(program *ast.Program, session Session) (object.Object, Session, time.Duration) {
	env := session.(*EvalSession).env

	start := time.Now()
	result := evaluator.Eval(program, env)
	duration := time.Since(start)
	return result, session, duration
}

func runVM(program *ast.Program, session Session) (object.Object, Session, time.Duration) {
	comp := compiler.New()

	if err := comp.Compile(program); err != nil {
		return err, session, 0
	}

	machine := vm.New(comp.Bytecode())

	start := time.Now()
	err := machine.Run()
	duration := time.Since(start)
	if err != nil {
		return object.NewError(nil, err), session, 0
	}

	result := machine.LastPoppedStackElem()
	return result, session, duration
}
