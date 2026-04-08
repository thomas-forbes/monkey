package runner

import (
	"monkey/ast"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/object"
	"monkey/vm"
	"time"
)

type Session interface {
	ExecProgram(program *ast.Program) (object.Object, time.Duration)
}

func NewSession(engine Engine) (Session, *object.Error) {
	var s Session
	switch engine {
	case INTERPRETER:
		s = NewEvalSession()
	case VM:
		s = NewVMSession()
	default:
		panic("unknown engine: " + string(engine))
	}

	stdProgram, err := ParseCode(STD)
	if err != nil {
		return nil, err
	}

	s.ExecProgram(stdProgram)
	return s, nil
}

// EVAL

type EvalSession struct {
	env *object.Environment
}

func NewEvalSession() *EvalSession {
	s := &EvalSession{env: object.NewEnvironment()}
	return s
}

func (s *EvalSession) ExecProgram(program *ast.Program) (object.Object, time.Duration) {
	env := s.env

	start := time.Now()
	result := evaluator.Eval(program, env)
	duration := time.Since(start)

	if result == nil {
		result = object.NULL
	}

	return result, duration
}

// VM

type VMSession struct {
	symbolTable *compiler.SymbolTable
	constants   []object.Object
	stack       []object.Object
	sp          int
}

func NewVMSession() *VMSession {
	s := &VMSession{
		symbolTable: compiler.NewMasterSymbolTable(),
		constants:   []object.Object{},
		stack:       make([]object.Object, vm.StackSize),
		sp:          0,
	}
	return s
}

func (s *VMSession) ExecProgram(program *ast.Program) (object.Object, time.Duration) {
	comp := compiler.NewWithState(s.symbolTable, s.constants)

	if err := comp.Compile(program); err != nil {
		return err, 0
	}

	machine := vm.NewWithState(comp.Bytecode(), s.stack, s.sp)

	start := time.Now()
	err := machine.Run()
	duration := time.Since(start)
	if err != nil {
		return object.NewError(nil, err), 0
	}

	s.sp = machine.GetSP()

	result := machine.LastPoppedStackElem()
	return result, duration

}
