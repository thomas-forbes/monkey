package runner

import (
	"fmt"
	"monkey/object"
	"testing"
)

type specCase struct {
	name     string
	input    string
	expected interface{}
}

type specGroup struct {
	name  string
	cases []specCase
}

func TestSpec(t *testing.T) {
	groups := []specGroup{
		{
			name: "arithmetic",
			cases: []specCase{
				{name: "integer literal", input: "5", expected: 5},
				{name: "integer literal negative", input: "-10", expected: -10},
				{name: "simple add", input: "1 + 2", expected: 3},
				{name: "simple subtract", input: "1 - 2", expected: -1},
				{name: "simple multiply", input: "1 * 2", expected: 2},
				{name: "simple divide", input: "4 / 2", expected: 2},
				{name: "chained arithmetic", input: "5 + 5 + 5 + 5 - 10", expected: 10},
				{name: "mixed precedence", input: "5 * 2 + 10", expected: 20},
				{name: "precedence reversed", input: "5 + 2 * 10", expected: 25},
				{name: "grouped expression", input: "2 * (5 + 10)", expected: 30},
				{name: "negative arithmetic", input: "-50 + 100 + -50", expected: 0},
				{name: "operator precedence", input: "(5 + 10 * 2 + 15 / 3) * 2 + -10", expected: 50},
			},
		},
		{
			name: "booleans",
			cases: []specCase{
				{name: "true literal", input: "true", expected: true},
				{name: "false literal", input: "false", expected: false},
				{name: "less than", input: "1 < 2", expected: true},
				{name: "greater than", input: "1 > 2", expected: false},
				{name: "integer equality", input: "1 == 1", expected: true},
				{name: "integer inequality", input: "1 != 2", expected: true},
				{name: "boolean equality", input: "true == false", expected: false},
				{name: "comparison to boolean", input: "(1 < 2) == true", expected: true},
				{name: "bang true", input: "!true", expected: false},
				{name: "bang false", input: "!false", expected: true},
				{name: "bang integer", input: "!5", expected: false},
				{name: "double bang integer", input: "!!5", expected: true},
				{name: "bang nullish if", input: "!(if (false) { 5; })", expected: true},
			},
		},
		{
			name: "conditionals",
			cases: []specCase{
				{name: "if true branch", input: "if (true) { 10 }", expected: 10},
				{name: "if else true branch", input: "if (true) { 10 } else { 20 }", expected: 10},
				{name: "if else false branch", input: "if (false) { 10 } else { 20 }", expected: 20},
				{name: "if integer truthy", input: "if (1) { 10 }", expected: 10},
				{name: "if null branch", input: "if (false) { 10 }", expected: nil},
				{name: "else if branch", input: "if (1 > 2) { 10 } else if (1 < 2) { 20 } else { 30 }", expected: 20},
				{name: "final else branch", input: "if (1 > 2) { 10 } else if (1 > 2) { 20 } else { 30 }", expected: 30},
				{name: "nested if falsey", input: "if ((if (false) { 10 })) { 10 } else { 20 }", expected: 20},
			},
		},
		{
			name: "returns",
			cases: []specCase{
				{name: "top level return", input: "return 10;", expected: 10},
				{name: "return skips following expression", input: "return 10; 9;", expected: 10},
				{name: "return skips following statement", input: "9; return 2 * 5; 9;", expected: 10},
				{name: "nested return statement", input: "if (10 > 1) { if (10 > 1) { return 10; } return 1; }", expected: 10},
				{name: "if else return branch", input: "if (1 > 2) { return 10; } else if (1 < 2) { return 20; } else { return 30; }", expected: 20},
			},
		},
		{
			name: "bindings",
			cases: []specCase{
				{name: "simple let", input: "let a = 5; a;", expected: 5},
				{name: "let expression value", input: "let a = 5 * 5; a;", expected: 25},
				{name: "let chain", input: "let a = 5; let b = a; let c = a + b + 5; c;", expected: 15},
				{name: "mutable assignment", input: "let mut a = 5; a = 3; a;", expected: 3},
				{name: "assignment preserves prior binding", input: "let mut a = 5; let b = a; a = 3; b;", expected: 5},
				{name: "underscore discard", input: "let _ = 5; _;", expected: nil},
				{name: "mutable function parameters", input: "let mut sum = fn(mut a, mut b) { a = a + 1; b = b + 1; a + b }; sum(1, 2);", expected: 5},
				{name: "immutable function parameters", input: "let mut sum = fn(a, b) { a = a + 1; b = b + 1; a + b }; sum(1, 2);", expected: errorObject("cannot assign to immutable variable: a")},
			},
		},
		{
			name: "functions",
			cases: []specCase{
				{name: "identity function", input: "let identity = fn(x) { x; }; identity(5);", expected: 5},
				{name: "identity return function", input: "let identity = fn(x) { return x; }; identity(5);", expected: 5},
				{name: "double function", input: "let double = fn(x) { x * 2; }; double(5);", expected: 10},
				{name: "add function", input: "let add = fn(x, y) { x + y; }; add(5, 5);", expected: 10},
				{name: "nested function application", input: "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", expected: 20},
				{name: "anonymous function", input: "fn(x) { x; }(5)", expected: 5},
				{name: "function without args", input: "fn() { 5 + 10; }();", expected: 15},
				{name: "function without return value", input: "fn() { }();", expected: nil},
				{name: "function return statement", input: "fn() { return 99; return 100; }();", expected: 99},
				{name: "global let in function", input: "let one = 1; let two = 2; one + two", expected: 3},
				{name: "let statement scopes", input: "fn() { let num = 1; num }();", expected: 1},
				{name: "first class function", input: "let returnsOne = fn() { 1; }; let returnsOneReturner = fn() { returnsOne; }; returnsOneReturner()();", expected: 1},
				{name: "function with arguments and globals", input: "let globalNum = 10; let sum = fn(a, b) { let c = a + b; c + globalNum; }; let outer = fn() { sum(1, 2) + sum(3, 4) + globalNum; }; outer() + globalNum;", expected: 50},
			},
		},
		{
			name: "closures",
			cases: []specCase{
				{name: "simple closure", input: "let newAdder = fn(x) { fn(y) { x + y }; }; let addTwo = newAdder(2); addTwo(2);", expected: 4},
				{name: "captured argument", input: "let newClosure = fn(a) { fn() { a; }; }; let closure = newClosure(99); closure();", expected: 99},
				{name: "captured sum", input: "let newAdder = fn(a, b) { let c = a + b; fn(d) { c + d }; }; let adder = newAdder(1, 2); adder(8);", expected: 11},
				{name: "nested closure", input: "let newAdderOuter = fn(a, b) { let c = a + b; fn(d) { let e = d + c; fn(f) { e + f; }; }; }; let newAdderInner = newAdderOuter(1, 2); let adder = newAdderInner(3); adder(8);", expected: 14},
				{name: "multiple enclosed fns", input: "let newClosure = fn(a, b) { let one = fn() { a; }; let two = fn() { b; }; fn() { one() + two(); }; }; let closure = newClosure(9, 90); closure();", expected: 99},
			},
		},
		{
			name: "strings",
			cases: []specCase{
				{name: "string literal", input: `"Hello World!"`, expected: "Hello World!"},
				{name: "simple string", input: `"monkey"`, expected: "monkey"},
				{name: "string concat", input: `"Hello" + " " + "World!"`, expected: "Hello World!"},
				{name: "multi concat", input: `"mon" + "key" + "banana"`, expected: "monkeybanana"},
			},
		},
		{
			name: "builtins",
			cases: []specCase{
				{name: "len empty string", input: `len("")`, expected: 0},
				{name: "len string", input: `len("four")`, expected: 4},
				{name: "len long string", input: `len("hello world")`, expected: 11},
				{name: "len array", input: `len([1, 2, 3])`, expected: 3},
				{name: "len empty array", input: `len([])`, expected: 0},
				{name: "append array", input: `append([], 1)`, expected: []int{1}},
				{name: "puts builtin", input: `puts("hello", "world!")`, expected: nil},
				{name: "len wrong type", input: `len(1)`, expected: errorObject("argument to `len` not supported, got INTEGER")},
				{name: "len wrong arg count", input: `len("one", "two")`, expected: errorObject("wrong number of arguments. got=2, want=1")},
				{name: "append wrong type", input: `append(1, 1)`, expected: errorObject("argument to `append` must be ARRAY, got INTEGER")},
				{name: "function wrong arg count", input: `fn(a, b) { a + b; }(1);`, expected: errorObject("wrong number of arguments: want=2, got=1")},
			},
		},
		{
			name: "arrays",
			cases: []specCase{
				{name: "empty array", input: "[]", expected: []int{}},
				{name: "array literal", input: "[1, 2 * 2, 3 + 3]", expected: []int{1, 4, 6}},
				{name: "array simple literal", input: "[1, 2, 3]", expected: []int{1, 2, 3}},
				{name: "array index 0", input: "[1, 2, 3][0]", expected: 1},
				{name: "array index 1", input: "[1, 2, 3][1]", expected: 2},
				{name: "array nested index", input: "[[1, 1, 1]][0][0]", expected: 1},
				{name: "array expression index", input: "[1, 2, 3][0 + 2]", expected: 3},
				{name: "array negative index", input: "[1, 2, 3][-1]", expected: 3},
				{name: "array negative index 2", input: "[1, 2, 3][-2]", expected: 2},
				{name: "array missing index", input: "[1, 2, 3][99]", expected: nil},
				{name: "empty array missing index", input: "[][0]", expected: nil},
			},
		},
		{
			name: "hashes",
			cases: []specCase{
				{name: "empty hash", input: "{}", expected: map[object.HashKey]int64{}},
				{name: "integer hash", input: "{1: 2, 2: 3}", expected: map[object.HashKey]int64{
					(&object.Integer{Value: 1}).HashKey(): 2,
					(&object.Integer{Value: 2}).HashKey(): 3,
				}},
				{name: "expression hash", input: "{1 + 1: 2 * 2, 3 + 3: 4 * 4}", expected: map[object.HashKey]int64{
					(&object.Integer{Value: 2}).HashKey(): 4,
					(&object.Integer{Value: 6}).HashKey(): 16,
				}},
				{name: "mixed hash literal", input: `let two = "two"; {"one": 10 - 9, two: 1 + 1, "thr" + "ee": 6 / 2, 4: 4, true: 5, false: 6}`, expected: map[object.HashKey]int64{
					(&object.String{Value: "one"}).HashKey():   1,
					(&object.String{Value: "two"}).HashKey():   2,
					(&object.String{Value: "three"}).HashKey(): 3,
					(&object.Integer{Value: 4}).HashKey():      4,
					object.TRUE.HashKey():                      5,
					object.FALSE.HashKey():                     6,
				}},
				{name: "hash string index", input: `{"foo": 5}["foo"]`, expected: 5},
				{name: "hash integer index", input: `{1: 1, 2: 2}[2]`, expected: 2},
				{name: "hash bool index", input: `{true: 5}[true]`, expected: 5},
				{name: "hash missing key", input: `{}["foo"]`, expected: nil},
				{name: "hash missing integer key", input: `{1: 1}[0]`, expected: nil},
				{name: "unusable hash key", input: `{"name": "Monkey"}[fn(x) { x }];`, expected: errorObject("unusable as hash key: FUNCTION")},
			},
		},
		{
			name: "errors",
			cases: []specCase{
				{name: "type mismatch", input: "5 + true;", expected: errorObject("type mismatch: INTEGER + BOOLEAN")},
				{name: "type mismatch after statement", input: "5 + true; 5;", expected: errorObject("type mismatch: INTEGER + BOOLEAN")},
				{name: "minus boolean", input: "-true", expected: errorObject("unknown operator: -BOOLEAN")},
				{name: "boolean arithmetic", input: "true + false;", expected: errorObject("unknown operator: BOOLEAN + BOOLEAN")},
				{name: "nested boolean arithmetic", input: "if (10 > 1) { true + false; }", expected: errorObject("unknown operator: BOOLEAN + BOOLEAN")},
				{name: "unknown identifier", input: "foobar", expected: errorObject("identifier not found: foobar")},
				{name: "unknown string operator", input: `"Hello" - "World"`, expected: errorObject("unknown operator: STRING - STRING")},
				{name: "reinitialize variable", input: `let a = 5; let a = 3;`, expected: errorObject("cannot reinitialize variable: a")},
				{name: "reassign immutable variable", input: `let a = 5; a = 3;`, expected: errorObject("cannot reassign unmutable variable: a")},
			},
		},
		{
			name: "loops",
			cases: []specCase{
				{name: "conditional for", input: `let mut i = 0; for i < 5 { i = i + 1; } i;`, expected: 5},
				{name: "conditional for continue", input: `let mut i = 0; let mut sum = 0; for i < 5 { i = i + 1; if (i == 3) { continue; } sum = sum + i; } sum;`, expected: 12},
				{name: "conditional for break", input: `let mut i = 0; let mut sum = 0; for true { if (i == 3) { break; } sum = sum + i; i = i + 1; } sum;`, expected: 3},
				{name: "conditional for return", input: `let mut i = 0; for i < 5 { if (i == 2) { return i; } i = i + 1; }`, expected: 2},
				{name: "conditional for break value", input: `let mut i = 0; for i < 5 { if (i == 2) { break 99; } i = i + 1; }`, expected: 99},
				{name: "for in array values", input: `let a = [1, 2, 3]; let mut sum = 0; for _, x in a { sum = sum + x } sum;`, expected: 6},
				{name: "for in array indexed", input: `let a = [1, 2, 3]; let mut sum = 0; for i, x in a { sum = sum + x * i } sum;`, expected: 8},
				{name: "for in hash keys", input: `let a = {1: 2, 3: 4, 5: 6}; let mut sum = 0; for k in a { sum = sum + k } sum;`, expected: 9},
				{name: "for in hash values", input: `let a = {1: 2, 3: 4, 5: 6}; let mut sum = 0; for _, v in a { sum = sum + v } sum;`, expected: 12},
				{name: "for in range", input: `let mut sum = 0; for i in 0..5 { sum = sum + i } sum;`, expected: 10},
				{name: "for in negative range", input: `let mut sum = 0; for i in -5..0 { sum = sum + i } sum;`, expected: -15},
			},
		},
		{
			name: "recursion",
			cases: []specCase{
				{name: "recursive function", input: `let countDown = fn(x) { if (x == 0) { return 0; } else { countDown(x - 1); } }; countDown(1);`, expected: 0},
				{name: "recursive function in wrapper", input: `let countDown = fn(x) { if (x == 0) { return 0; } else { countDown(x - 1); } }; let wrapper = fn() { countDown(1); }; wrapper();`, expected: 0},
				{name: "recursive fibonacci", input: `let fibonacci = fn(x) { if (x == 0) { return 0; } else { if (x == 1) { return 1; } else { fibonacci(x - 1) + fibonacci(x - 2); } } }; fibonacci(15);`, expected: 610},
			},
		},
	}

	engines := []Engine{INTERPRETER, VM}
	for _, group := range groups {
		t.Run(group.name, func(t *testing.T) {
			for _, tt := range group.cases {
				for _, engine := range engines {
					t.Run(fmt.Sprintf("%s/%s", engine, tt.name), func(t *testing.T) {
						result, _, _ := RunProgram(engine, tt.input, NewEnvironment(engine))
						assertExpectedObject(t, tt.expected, result)
					})
				}
			}
		})
	}
}

func errorObject(message string) *object.Error {
	return &object.Error{Message: message}
}

func assertExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	if expected == nil {
		if actual != object.NULL {
			t.Fatalf("object is not NULL. got=%T (%+v)", actual, actual)
		}
		return
	}

	switch expected := expected.(type) {
	case int:
		integer, ok := actual.(*object.Integer)
		if !ok {
			t.Fatalf("object is not Integer. got=%T (%+v)", actual, actual)
		}
		if integer.Value != int64(expected) {
			t.Fatalf("integer has wrong value. got=%d, want=%d", integer.Value, expected)
		}
	case bool:
		boolean, ok := actual.(*object.Boolean)
		if !ok {
			t.Fatalf("object is not Boolean. got=%T (%+v)", actual, actual)
		}
		if boolean.Value != expected {
			t.Fatalf("boolean has wrong value. got=%t, want=%t", boolean.Value, expected)
		}
	case string:
		str, ok := actual.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", actual, actual)
		}
		if str.Value != expected {
			t.Fatalf("string has wrong value. got=%q, want=%q", str.Value, expected)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", actual, actual)
		}
		if len(array.Elements) != len(expected) {
			t.Fatalf("array has wrong num of elements. got=%d, want=%d", len(array.Elements), len(expected))
		}
		for i, expectedElem := range expected {
			integer, ok := array.Elements[i].(*object.Integer)
			if !ok {
				t.Fatalf("array element %d is not Integer. got=%T (%+v)", i, array.Elements[i], array.Elements[i])
			}
			if integer.Value != int64(expectedElem) {
				t.Fatalf("array element %d has wrong value. got=%d, want=%d", i, integer.Value, expectedElem)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Fatalf("object is not Hash. got=%T (%+v)", actual, actual)
		}
		if len(hash.Pairs) != len(expected) {
			t.Fatalf("hash has wrong number of pairs. got=%d, want=%d", len(hash.Pairs), len(expected))
		}
		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Fatalf("missing hash pair for key %+v", expectedKey)
			}
			integer, ok := pair.Value.(*object.Integer)
			if !ok {
				t.Fatalf("hash value is not Integer. got=%T (%+v)", pair.Value, pair.Value)
			}
			if integer.Value != expectedValue {
				t.Fatalf("hash value has wrong value. got=%d, want=%d", integer.Value, expectedValue)
			}
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Fatalf("object is not Error. got=%T (%+v)", actual, actual)
		}
		if errObj.Message != expected.Message {
			t.Fatalf("error has wrong message. got=%q, want=%q", errObj.Message, expected.Message)
		}
	default:
		t.Fatalf("unhandled expected type %T", expected)
	}
}
