package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

const fibProgram = `
let fib = fn(x) {
	if (x < 2) {
		x
	} else {
		fib(x - 1) + fib(x - 2)
	}
};

return fib(10);
`

func BenchmarkFibEval(b *testing.B) {
	l := lexer.New(fibProgram)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		b.Fatalf("parser errors: %v", p.Errors())
	}

	result := Eval(program, object.NewEnvironment())
	intResult, ok := result.(*object.Integer)
	if !ok {
		b.Fatalf("unexpected result type: %T", result)
	}
	if intResult.Value != 55 {
		b.Fatalf("unexpected result value: got=%d want=55", intResult.Value)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Eval(program, object.NewEnvironment())
	}
}
