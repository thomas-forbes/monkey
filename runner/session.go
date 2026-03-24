package runner

import (
	"monkey/compiler"
	"monkey/object"
	"monkey/vm"
)

type Session interface {
	session()
}

func NewSession(engine Engine) Session {
	switch engine {
	case INTERPRETER:
		return &EvalSession{env: object.NewEnvironment()}
	case VM:
		sp := 0
		return &VMSession{
			symbolTable: compiler.NewMasterSymbolTable(),
			constants:   []object.Object{},
			stack:       make([]object.Object, vm.StackSize),
			sp:          &sp,
		}
	default:
		panic("unknown engine: " + string(engine))
	}
}

type EvalSession struct {
	env *object.Environment
}

func (s *EvalSession) session() {}

type VMSession struct {
	symbolTable *compiler.SymbolTable
	constants   []object.Object
	stack       []object.Object
	sp          *int
}

func (s *VMSession) session() {}
