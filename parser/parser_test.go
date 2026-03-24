package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
		expectedMutable    bool
	}{
		{"let x = 5;", "x", 5, false},
		{"let mut y = true;", "y", true, true},
		{"let foobar = y;", "foobar", "y", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
		mutable := stmt.(*ast.LetStatement).Initialization.Mutable
		if mutable != tt.expectedMutable {
			t.Errorf("stmt.Mutable not %t. got=%t", tt.expectedMutable, mutable)
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Initialization.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Initialization.Name.Value)
		return false
	}
	if letStmt.Initialization.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Initialization.Name.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			return
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	if !testLiteralExpression(t, stmt.Expression, "foobar") {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true; false;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 2 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt1, ok1 := program.Statements[0].(*ast.ExpressionStatement)
	stmt2, ok2 := program.Statements[1].(*ast.ExpressionStatement)
	if !ok1 || !ok2 {
		t.Fatalf("program.Statements[0] or [1] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	if !testLiteralExpression(t, stmt1.Expression, true) {
		return
	}
	if !testLiteralExpression(t, stmt2.Expression, false) {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))

		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestAssignmentExpression(t *testing.T) {
	input := "x = 5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.AssignmentExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Name, "x") {
		return
	}
	if !testLiteralExpression(t, exp.Value, 5) {
		return
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"a = [1, 2, 3, 4][b * c] * d",
			"a = (([1, 2, 3, 4][(b * c)]) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)

	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		branches []struct {
			condition string
			body      string
		}
	}{
		{
			"if (x < y) { x }",
			[]struct{ condition, body string }{
				{"(x < y)", "x"},
			},
		},
		{
			"if (x < y) { x } else { y }",
			[]struct{ condition, body string }{
				{"(x < y)", "x"},
				{"", "y"},
			},
		},
		{
			"if (x < y) { x } else if (y > x) { y } else { z }",
			[]struct{ condition, body string }{
				{"(x < y)", "x"},
				{"(y > x)", "y"},
				{"", "z"},
			},
		},
		{
			"if (x < y) { return x; } else if (y > x) { return y; } else { return z; }",
			[]struct{ condition, body string }{
				{"(x < y)", "return x;"},
				{"(y > x)", "return y;"},
				{"", "return z;"},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt is not ast.IfExpression. got=%T", stmt.Expression)
		}

		if len(exp.Branches) != len(tt.branches) {
			t.Fatalf("exp.Branches does not contain %d branches. got=%d\n",
				len(tt.branches), len(exp.Branches))
		}

		for i, branch := range tt.branches {
			expBranch := exp.Branches[i]
			if expBranch.Condition == nil {
				if branch.condition != "" {
					t.Errorf("exp.Branches[%d].Condition is not %s. got=nil",
						i, branch.condition)
				}
			} else if expBranch.Condition.String() != branch.condition {
				t.Errorf("exp.Branches[%d].Condition is not %s. got=%s",
					i, branch.condition, exp.Branches[i].Condition.String())
			}

			if exp.Branches[i].Body.String() != branch.body {
				t.Errorf("exp.Branches[%d].Body is not %s. got=%s",
					i, branch.body, exp.Branches[i].Body.String())
			}
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0].Name, "x")
	testLiteralExpression(t, function.Parameters[1].Name, "y")
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))

	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i].Name, ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
		if len(array.Elements) != 3 {
			t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
		}
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
		if len(hash.Pairs) != 3 {
			t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
		}
	}
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}
	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func TestMultiStatementProgram(t *testing.T) {
	tests := []struct {
		input string
		stmts []ast.Statement
	}{
		{
			"let x = 5; let y = 10;",
			[]ast.Statement{&ast.LetStatement{}, &ast.LetStatement{}},
		},
		{
			"let add = fn(x, y) { x + y }; let x = add(5, 10);",
			[]ast.Statement{&ast.LetStatement{}, &ast.LetStatement{}},
		},
		{
			"let check = fn(x, y) { if (x > y) { return x; } else if (y > x) { return y; } else { return null; } }; let x = add(5, 10);",
			[]ast.Statement{&ast.LetStatement{}, &ast.LetStatement{}},
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != len(tt.stmts) {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", len(tt.stmts), len(program.Statements))
		}
		for i, testStmt := range tt.stmts {
			realStmt := program.Statements[i]
			switch testStmt.(type) {
			case *ast.LetStatement:
				if _, ok := realStmt.(*ast.LetStatement); !ok {
					t.Fatalf("statement %d is wrong type: got=%T, want=*ast.LetStatement", i, realStmt)
				}
			case *ast.ReturnStatement:
				if _, ok := realStmt.(*ast.ReturnStatement); !ok {
					t.Fatalf("statement %d is wrong type: got=%T, want=*ast.ReturnStatement", i, realStmt)
				}
			case *ast.ExpressionStatement:
				if _, ok := realStmt.(*ast.ExpressionStatement); !ok {
					t.Fatalf("statement %d is wrong type: got=%T, want=*ast.ExpressionStatement", i, realStmt)
				}
			default:
				t.Fatalf("unsupported expected statement type %T", testStmt)
			}
		}
	}
}

func TestForRangeStatement(t *testing.T) {
	tests := []struct {
		input           string
		indexIdentifier string
		valueIdentifier string
		rangeString     string
		bodyString      string
	}{
		{
			"for i, x in myArray { x }",
			"i",
			"x",
			"myArray",
			"x",
		},
		{
			"for _, x in range(5) { return x; }",
			"_",
			"x",
			"range(5)",
			"return x;",
		},
		{
			"for i in range(4 + 4) { x = x + i; }",
			"i",
			"",
			"range((4 + 4))",
			"x = (x + i)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
		}

		clause, ok := stmt.Clause.(*ast.ForInClause)
		if !ok {
			t.Fatalf("stmt.Clause is not ast.ForRangeClause. got=%T", stmt.Clause)
		}

		if !testIdentifier(t, clause.Key, tt.indexIdentifier) {
			return
		}

		if len(tt.valueIdentifier) > 0 && !testIdentifier(t, clause.Value, tt.valueIdentifier) {
			return
		} else if len(tt.valueIdentifier) == 0 && clause.Value != nil {
			t.Errorf("stmt.Value should be nil. got=%T", clause.Value)
		}

		if clause.Iterable.String() != tt.rangeString {
			t.Errorf("stmt.Range.String() not %q. got=%q", tt.rangeString, clause.Iterable.String())
		}

		if stmt.Body.String() != tt.bodyString {
			t.Errorf("stmt.Body.String() not %q. got=%q", tt.bodyString, stmt.Body.String())
		}
	}
}

func TestForConditionStatement(t *testing.T) {
	tests := []struct {
		input      string
		condString string
		bodyString string
	}{
		{
			"for i < 10 { i = i + 1; }",
			"(i < 10)",
			"i = (i + 1)",
		},
		{
			"for true { i = i + 1; }",
			"true",
			"i = (i + 1)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
		}

		clause, ok := stmt.Clause.(*ast.ForConditionalClause)
		if !ok {
			t.Fatalf("stmt.Clause is not ast.ForConditionClause. got=%T", stmt.Clause)
		}

		if clause.Condition.String() != tt.condString {
			t.Errorf("stmt.Condition.String() not %q. got=%q", tt.condString, clause.Condition.String())
		}

		if stmt.Body.String() != tt.bodyString {
			t.Errorf("stmt.Body.String() not %q. got=%q", tt.bodyString, stmt.Body.String())
		}
	}
}

func TestRangeExpression(t *testing.T) {
	tests := []struct {
		input       string
		leftString  string
		rightString string
	}{
		{
			"0..10",
			"0",
			"10",
		},
		{
			"-1..-10",
			"(-1)",
			"(-10)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.RangeExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.RangeExpression. got=%T", stmt.Expression)
		}

		if exp.Left.String() != tt.leftString {
			t.Errorf("exp.Start.String() not %q. got=%q", tt.leftString, exp.Left.String())
		}

		if exp.Right.String() != tt.rightString {
			t.Errorf("exp.End.String() not %q. got=%q", tt.rightString, exp.Right.String())
		}
	}
}

func TestBreakStatement(t *testing.T) {
	input := "for i in range(10) { if (i == 5) { break 2 + 2; } }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements does not contain %d statements. got=%d\n", 1, len(stmt.Body.Statements))
	}

	expressionStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Body.Statements[0])
	}
	ifExpression, ok := expressionStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0].Expression is not ast.IfExpression. got=%T", expressionStmt.Expression)
	}
	if len(ifExpression.Branches) != 1 {
		t.Fatalf("ifExpression.Branches does not contain %d branches. got=%d\n", 1, len(ifExpression.Branches))
	}
	ifBranch := ifExpression.Branches[0]
	breakStmt, ok := ifBranch.Body.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("ifBranch.Body.Statements[0] is not ast.BreakStatement. got=%T", ifBranch.Body.Statements[0])
	}

	_, ok = breakStmt.Value.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("breakStmt.BreakValue is not ast.InfixExpression. got=%T", breakStmt.Value)
	}
}

func TestContinueStatement(t *testing.T) {
	input := "for i in range(10) { if (i == 5) { continue; } }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("stmt.Body.Statements does not contain %d statements. got=%d\n", 1, len(stmt.Body.Statements))
	}

	expressionStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Body.Statements[0])
	}
	ifExpression, ok := expressionStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0].Expression is not ast.IfExpression. got=%T", expressionStmt.Expression)
	}
	if len(ifExpression.Branches) != 1 {
		t.Fatalf("ifExpression.Branches does not contain %d branches. got=%d\n", 1, len(ifExpression.Branches))
	}
	ifBranch := ifExpression.Branches[0]
	_, ok = ifBranch.Body.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("ifBranch.Body.Statements[0] is not ast.ContinueStatement. got=%T", ifBranch.Body.Statements[0])
	}
}

func TestFunctionLiteralWithName(t *testing.T) {
	input := `let myFunction = fn() { };`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}

	if function.Name != "myFunction" {
		t.Fatalf("function literal name wrong. want 'myFunction', got=%q\n",
			function.Name)
	}
}
