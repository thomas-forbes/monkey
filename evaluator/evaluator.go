package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)

	// expressions
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.RangeExpression:
		return evalRangeExpression(node, env)
	case *ast.BreakStatement:
		value := Eval(node.Value, env)
		return &object.Break{Value: value}
	case *ast.ContinueStatement:
		return &object.Continue{}
	}
	return object.NULL
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
		_, ok := result.(object.ControlFlowSignal)
		if ok {
			return result
		}
	}
	return result
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unsupported type for negation: %s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("unsupported types for binary operation: %s %s", left.Type(), right.Type())
	default:
		return newError("unsupported types for binary operation: %s %s", left.Type(), right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s", left.Type(), right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	for _, branch := range ie.Branches {
		if branch.Condition == nil {
			return Eval(branch.Body, env)
		}

		condition := Eval(branch.Condition, env)
		if IsError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(branch.Body, env)
		}
	}
	return object.NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if entity, ok := env.Get(node.Value); ok {
		return entity.Object
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv, err := extendFunctionEnv(fn, args)
		if err != nil {
			return err
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnvironment(fn.Env)
	if len(args) != len(fn.Parameters) {
		return nil, &object.Error{Message: fmt.Sprintf("wrong number of arguments: want=%d, got=%d", len(fn.Parameters), len(args))}
	}
	for paramIdx, param := range fn.Parameters {
		_, ok := env.Set(param.Name.Value, args[paramIdx], param.Mutable, true)
		if !ok {
			return nil, &object.Error{Message: fmt.Sprintf("cannot reinitialize variable: %s", param.Name.Value)}
		}
	}
	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 {
		idx = idx + max + 1
	}
	if idx < 0 || idx > max {
		return object.NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if IsError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if IsError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return object.NULL
	}
	return pair.Value
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if IsError(val) {
		return val
	}
	name := node.Initialization.Name.Value
	if entity, ok := env.Get(name); ok && entity != object.NULL_ENTITY {
		return newError("cannot reinitialize variable: %s", name)
	}
	env.Set(name, val, node.Initialization.Mutable, true)
	return object.NULL
}

func evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if IsError(val) {
		return val
	}
	name := node.Name.Value
	entity, ok := env.Get(name)
	if !ok {
		return newError("identifier not found: %s", name)
	}
	if !entity.Mutable {
		return newError("cannot assign to immutable variable: %s", name)
	}
	env.Set(name, val, entity.Mutable, false)
	return val
}

func evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	switch node.Clause.(type) {
	case *ast.ForInClause:
		return evalForInStatement(node, node.Clause.(*ast.ForInClause), env)
	case *ast.ForConditionalClause:
		return evalForConditionalStatement(node, node.Clause.(*ast.ForConditionalClause), env)
	default:
		return newError("unknown for control clause: %T", node.Clause)
	}
}

func evalLoopBody(body *ast.BlockStatement, env *object.Environment) (object.Object, bool) {
	result := Eval(body, env)
	control, ok := result.(object.ControlFlowSignal)
	if !ok {
		return result, false
	}

	switch control.(type) {
	case *object.Break:
		return result.(*object.Break).Value, true
	case *object.Continue:
		return result, false
	case *object.ReturnValue:
		return result, true
	case *object.Error:
		return result, true
	default:
		return result, true
	}
}

func evalForInStatement(node *ast.ForStatement, clause *ast.ForInClause, env *object.Environment) object.Object {
	keyName := clause.Key.Value
	var valueName *string = nil
	if clause.Value != nil {
		valueName = &clause.Value.Value
	}

	switch iter := Eval(clause.Iterable, env).(type) {
	case *object.Array:
		for index, value := range iter.Elements {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(keyName, &object.Integer{Value: int64(index)}, false, true)
			if valueName != nil {
				loopEnv.Set(*valueName, value, false, true)
			}

			result, done := evalLoopBody(node.Body, loopEnv)
			if done {
				return result
			}
		}
	case *object.Hash:
		for _, pair := range iter.Pairs {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(keyName, pair.Key, false, true)
			if valueName != nil {
				loopEnv.Set(*valueName, pair.Value, false, true)
			}

			result, done := evalLoopBody(node.Body, loopEnv)
			if done {
				return result
			}
		}
	case *object.Range:
		if valueName != nil {
			return newError("cannot assign value in for-in loop without array range")
		}

		var increment int64 = 1
		if iter.Left > iter.Right {
			increment = -1
		}

		for index := iter.Left; (increment < 0 && index > iter.Right) || (increment > 0 && index < iter.Right); index = index + increment {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(keyName, &object.Integer{Value: int64(index)}, false, true)

			result, done := evalLoopBody(node.Body, loopEnv)
			if done {
				return result
			}
		}
	}

	return object.NULL
}

func evalForConditionalStatement(node *ast.ForStatement, clause *ast.ForConditionalClause, env *object.Environment) object.Object {
	for {
		condition := Eval(clause.Condition, env)
		if IsError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			return object.NULL
		}

		result, done := evalLoopBody(node.Body, env)
		if done {
			return result
		}
	}
}

func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if IsError(left) {
		return left
	}
	leftInt, ok := left.(*object.Integer)
	if !ok {
		return newError("left side of range expression must be an integer, got %s", left.Type())
	}
	right := Eval(node.Right, env)
	if IsError(right) {
		return right
	}
	rightInt, ok := right.(*object.Integer)
	if !ok {
		return newError("right side of range expression must be an integer, got %s", right.Type())
	}
	return &object.Range{Right: rightInt.Value, Left: leftInt.Value}
}
