package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: LocalScope, Index: 0},
		"b": {Name: "b", Scope: LocalScope, Index: 1, Mutable: true},
		"c": {Name: "c", Scope: LocalScope, Index: 0},
		"d": {Name: "d", Scope: LocalScope, Index: 1},
		"e": {Name: "e", Scope: LocalScope, Index: 0},
		"f": {Name: "f", Scope: LocalScope, Index: 1},
	}
	global := NewMasterSymbolTable()
	a, ok := global.Define("a", false)
	if !ok {
		t.Fatalf("expected define(a) to succeed")
	}
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	b, ok := global.Define("b", true)
	if !ok {
		t.Fatalf("expected define(b) to succeed")
	}
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}
	firstLocal := NewEnclosedSymbolTable(global)
	c, ok := firstLocal.Define("c", false)
	if !ok {
		t.Fatalf("expected define(c) to succeed")
	}
	if c != expected["c"] {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}
	d, ok := firstLocal.Define("d", false)
	if !ok {
		t.Fatalf("expected define(d) to succeed")
	}
	if d != expected["d"] {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}
	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e, ok := secondLocal.Define("e", false)
	if !ok {
		t.Fatalf("expected define(e) to succeed")
	}
	if e != expected["e"] {
		t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
	}
	f, ok := secondLocal.Define("f", false)
	if !ok {
		t.Fatalf("expected define(f) to succeed")
	}
	if f != expected["f"] {
		t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
	}
}
func TestResolveGlobal(t *testing.T) {
	global := NewMasterSymbolTable()
	global.Define("a", false)
	global.Define("b", false)
	expected := []Symbol{
		{Name: "a", Scope: LocalScope, Index: 0},
		{Name: "b", Scope: LocalScope, Index: 1},
	}
	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}
}

func TestResolve(t *testing.T) {
	global := NewMasterSymbolTable()
	global.Define("a", false)
	global.Define("b", false)

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c", false)
	firstLocal.Define("d", false)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)
	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: FreeScope, Index: 0},
				{Name: "b", Scope: FreeScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: FreeScope, Index: 0},
				{Name: "b", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}
	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, result)
			}
		}
	}
}

func TestResolveFreeSymbols(t *testing.T) {
	global := NewMasterSymbolTable()
	global.Define("a", false)
	global.Define("b", false)

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c", false)
	firstLocal.Define("d", false)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

	expected := []Symbol{
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "d", Scope: FreeScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Fatalf("name %s not resolvable", sym.Name)
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}

	if len(secondLocal.FreeSymbols) != 2 {
		t.Fatalf("wrong number of free symbols. got=%d", len(secondLocal.FreeSymbols))
	}

	wantFree := []Symbol{
		{Name: "c", Scope: LocalScope, Index: 0},
		{Name: "d", Scope: LocalScope, Index: 1},
	}
	for i, sym := range wantFree {
		if secondLocal.FreeSymbols[i] != sym {
			t.Errorf("wrong free symbol at %d. want=%+v, got=%+v", i, sym, secondLocal.FreeSymbols[i])
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewMasterSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		Symbol{Name: "a", Scope: BuiltinScope, Index: 0},
		Symbol{Name: "c", Scope: BuiltinScope, Index: 1},
		Symbol{Name: "e", Scope: BuiltinScope, Index: 2},
		Symbol{Name: "f", Scope: BuiltinScope, Index: 3},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	global := NewMasterSymbolTable()
	global.DefineFunctionName("a")

	expected := Symbol{Name: "a", Scope: FunctionScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v",
			expected.Name, expected, result)
	}
}

func TestBlockScope(t *testing.T) {
	outer := NewMasterSymbolTable()
	outer.Define("a", false)
	outer.Define("b", false)

	block := NewBlockSymbolTable(outer)

	// resolving outer symbols from block returns LocalScope, not FreeScope
	result, ok := block.Resolve("a")
	if !ok {
		t.Fatalf("name a not resolvable")
	}
	if result != (Symbol{Name: "a", Scope: LocalScope, Index: 0}) {
		t.Errorf("expected a=%+v, got=%+v", Symbol{Name: "a", Scope: LocalScope, Index: 0}, result)
	}

	// block should not generate free symbols
	if len(block.FreeSymbols) != 0 {
		t.Errorf("expected no free symbols, got=%d", len(block.FreeSymbols))
	}

	// definitions in block continue index from outer
	c, ok := block.Define("c", false)
	if !ok {
		t.Fatalf("expected define(c) to succeed")
	}
	if c != (Symbol{Name: "c", Scope: LocalScope, Index: 2}) {
		t.Errorf("expected c=%+v, got=%+v", Symbol{Name: "c", Scope: LocalScope, Index: 2}, c)
	}
}

func TestBlockScopeShadowing(t *testing.T) {
	outer := NewMasterSymbolTable()
	outer.Define("x", false) // index 0

	block := NewBlockSymbolTable(outer)
	x, ok := block.Define("x", false) // shadows outer x, gets new index
	if !ok {
		t.Fatalf("expected define(x) to succeed in block")
	}
	if x.Index != 1 {
		t.Errorf("expected shadowed x to have index 1, got=%d", x.Index)
	}

	// resolving x in block returns the block-local one
	resolved, ok := block.Resolve("x")
	if !ok {
		t.Fatalf("name x not resolvable in block")
	}
	if resolved.Index != 1 {
		t.Errorf("expected block x to resolve to index 1, got=%d", resolved.Index)
	}

	// resolving x in outer still returns index 0
	resolved, ok = outer.Resolve("x")
	if !ok {
		t.Fatalf("name x not resolvable in outer")
	}
	if resolved.Index != 0 {
		t.Errorf("expected outer x to resolve to index 0, got=%d", resolved.Index)
	}
}

func TestBlockScopeHighWaterMark(t *testing.T) {
	outer := NewMasterSymbolTable()
	outer.Define("a", false) // index 0

	block := NewBlockSymbolTable(outer)
	block.Define("b", false) // index 1
	block.Define("c", false) // index 2

	// simulate leaveBlock
	if block.numDefinitions > outer.numDefinitions {
		outer.numDefinitions = block.numDefinitions
	}

	if outer.numDefinitions != 3 {
		t.Errorf("expected outer numDefinitions=3, got=%d", outer.numDefinitions)
	}

	// new definition in outer continues from high-water-mark
	d, ok := outer.Define("d", false)
	if !ok {
		t.Fatalf("expected define(d) to succeed")
	}
	if d.Index != 3 {
		t.Errorf("expected d to have index 3, got=%d", d.Index)
	}
}

func TestBlockScopeInFunction(t *testing.T) {
	global := NewMasterSymbolTable()
	global.Define("a", false)

	fn := NewEnclosedSymbolTable(global)
	fn.Define("b", false)

	block := NewBlockSymbolTable(fn)
	block.Define("c", false)

	// c is local to the function frame
	result, ok := block.Resolve("c")
	if !ok {
		t.Fatalf("name c not resolvable")
	}
	if result.Scope != LocalScope {
		t.Errorf("expected c to have LocalScope, got=%s", result.Scope)
	}

	// a is free (crosses function boundary)
	result, ok = block.Resolve("a")
	if !ok {
		t.Fatalf("name a not resolvable")
	}
	if result.Scope != FreeScope {
		t.Errorf("expected a to have FreeScope, got=%s", result.Scope)
	}

	// b is local (same function frame, crosses block boundary)
	result, ok = block.Resolve("b")
	if !ok {
		t.Fatalf("name b not resolvable")
	}
	if result.Scope != LocalScope {
		t.Errorf("expected b to have LocalScope, got=%s", result.Scope)
	}
}

func TestShadowingFunctionName(t *testing.T) {
	global := NewMasterSymbolTable()
	global.DefineFunctionName("a")
	_, ok := global.Define("a", false)
	if ok {
		t.Fatalf("expected define(a) to fail when function name already exists")
	}

	expected := Symbol{Name: "a", Scope: FunctionScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got=%+v",
			expected.Name, expected, result)
	}
}
