package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{Name: "a", Index: 0},
		"b": Symbol{Name: "b", Index: 1},
		"c": Symbol{Name: "c", Index: 0},
		"d": Symbol{Name: "d", Index: 1},
		"e": Symbol{Name: "e", Index: 0},
		"f": Symbol{Name: "f", Index: 1},
	}
	global := NewSymbolTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}
	firstLocal := NewEnclosedSymbolTable(global)
	c := firstLocal.Define("c")
	if c != expected["c"] {
		t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
	}
	d := firstLocal.Define("d")
	if d != expected["d"] {
		t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
	}
	secondLocal := NewEnclosedSymbolTable(firstLocal)
	e := secondLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
	}
	f := secondLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
	}
}
func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")
	expected := []Symbol{
		Symbol{Name: "a", Index: 0},
		Symbol{Name: "b", Index: 1},
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
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")
	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				Symbol{Name: "a", Index: 0},
				Symbol{Name: "b", Index: 1},
				Symbol{Name: "c", Index: 0},
				Symbol{Name: "d", Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				Symbol{Name: "a", Index: 0},
				Symbol{Name: "b", Index: 1},
				Symbol{Name: "e", Index: 0},
				Symbol{Name: "f", Index: 1},
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
