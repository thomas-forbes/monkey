package compiler

import "monkey/object"

type SymbolScope string

const (
	LocalScope    SymbolScope = "LOCAL"
	FreeScope     SymbolScope = "FREE"
	BuiltinScope  SymbolScope = "BUILTIN"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name    string
	Scope   SymbolScope
	Index   int
	Mutable bool
}

type SymbolTable struct {
	Outer       *SymbolTable
	FreeSymbols []Symbol

	store          map[string]Symbol
	numDefinitions int
}

func newSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: s, FreeSymbols: free, numDefinitions: 0}
}

func NewMasterSymbolTable() *SymbolTable {
	symbolTable := newSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}
	return symbolTable
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := newSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string, mutable bool) (Symbol, bool) {
	symbol := Symbol{Name: name, Scope: LocalScope, Index: s.numDefinitions, Mutable: mutable}
	if _, ok := s.store[name]; ok {
		return symbol, false
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol, true
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}
		if obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.defineFree(obj)
		return free, true
	}
	return obj, ok
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1, Mutable: original.Mutable}
	symbol.Scope = FreeScope

	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope, Mutable: false}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope, Mutable: false}
	s.store[name] = symbol
	return symbol
}
