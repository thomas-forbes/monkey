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
	"strings"
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

func NewEnvironment(engine Engine) *object.Environment {
	env := object.NewEnvironment()
	if strings.TrimSpace(stdlib) == "" {
		return env
	}

	result, loadedEnv, _ := RunProgram(engine, stdlib, env)
	if result.Type() == object.ERROR_OBJ {
		panic(fmt.Errorf("error loading stdlib: %s", result.Inspect()))
	}

	return loadedEnv
}

func RunProgram(engine Engine, input string, env *object.Environment) (object.Object, *object.Environment, time.Duration) {
	if env == nil {
		env = object.NewEnvironment()
	}

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return &object.Error{Message: strings.Join(p.Errors(), "\n")}, env, 0
	}

	var result object.Object
	var duration time.Duration

	switch engine {
	case INTERPRETER:
		start := time.Now()
		evalResult, err := runEvaluator(program, env)
		duration = time.Since(start)
		if err != nil {
			return &object.Error{Message: err.Error()}, env, duration
		}
		result = evalResult
	case VM:
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			return &object.Error{Message: err.Error()}, env, 0
		}
		// fmt.Println(comp.Bytecode().Instructions)
		machine := vm.New(comp.Bytecode())

		start := time.Now()
		err = runVM(machine)
		duration = time.Since(start)
		if err != nil {
			return &object.Error{Message: err.Error()}, env, duration
		}

		result = machine.LastPoppedStackElem()
	default:
		return &object.Error{Message: fmt.Sprintf("unknown engine: %s", engine)}, env, 0
	}

	if result == nil {
		return object.NULL, env, duration
	}

	return result, env, duration
}

func runVM(machine *vm.VM) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	return machine.Run()
}

func runEvaluator(program *ast.Program, env *object.Environment) (result object.Object, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	return evaluator.Eval(program, env), nil
}
