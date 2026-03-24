package runner

import (
	"monkey/compiler"
	"monkey/object"
)

type Session interface {
	session()
}

func NewSession(engine Engine) Session {
	switch engine {
	case INTERPRETER:
		return &EvalSession{env: object.NewEnvironment()}
	case VM:
		return &VMSession{
			symbolTable: compiler.NewMasterSymbolTable(),
			constants:   &[]object.Object{},
			// globals:     &[]object.Object{},
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
	constants   *[]object.Object
	// globals     *[]object.Object
}

func (s *VMSession) session() {}
