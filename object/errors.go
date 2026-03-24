package object

import (
	"fmt"
	"monkey/ast"
	"monkey/parser"
	"monkey/token"
	"strings"
)

type UnknownIdentifier struct {
	Name string
}

func (e UnknownIdentifier) Error() string {
	return fmt.Sprintf("identifier not found: %s", e.Name)
}

func (e UnknownIdentifier) String() string {
	return e.Error()
}

func NewUnknownIdentifier(tok *token.Token, name string) *Error {
	return NewError(tok, UnknownIdentifier{Name: name})
}

type CannotReinitializeVariable struct {
	Name string
}

func (e CannotReinitializeVariable) Error() string {
	return fmt.Sprintf("cannot reinitialize variable: %s", e.Name)
}

func (e CannotReinitializeVariable) String() string {
	return e.Error()
}

func NewCannotReinitializeVariable(tok *token.Token, name string) *Error {
	return NewError(tok, CannotReinitializeVariable{Name: name})
}

type CannotAssignImmutableVariable struct {
	Name string
}

func (e CannotAssignImmutableVariable) Error() string {
	return fmt.Sprintf("cannot assign to immutable variable: %s", e.Name)
}

func (e CannotAssignImmutableVariable) String() string {
	return e.Error()
}

func NewCannotAssignImmutableVariable(tok *token.Token, name string) *Error {
	return NewError(tok, CannotAssignImmutableVariable{Name: name})
}

type WrongArgumentCount struct {
	Want int
	Got  int
}

func (e WrongArgumentCount) Error() string {
	return fmt.Sprintf("wrong number of arguments: want=%d, got=%d", e.Want, e.Got)
}

func (e WrongArgumentCount) String() string {
	return e.Error()
}

func NewWrongArgumentCount(tok *token.Token, want int, got int) *Error {
	return NewError(tok, WrongArgumentCount{Want: want, Got: got})
}

type UnsupportedBinaryOperation struct {
	LeftType  string
	RightType string
}

func (e UnsupportedBinaryOperation) Error() string {
	return fmt.Sprintf("unsupported types for binary operation: %s %s", e.LeftType, e.RightType)
}

func (e UnsupportedBinaryOperation) String() string {
	return e.Error()
}

func NewUnsupportedBinaryOperation(tok *token.Token, leftType string, rightType string) *Error {
	return NewError(tok, UnsupportedBinaryOperation{LeftType: leftType, RightType: rightType})
}

type UnsupportedUnaryOperation struct {
	Operator    string
	OperandType string
}

func (e UnsupportedUnaryOperation) Error() string {
	switch e.Operator {
	case "-":
		return fmt.Sprintf("unsupported type for negation: %s", e.OperandType)
	default:
		return fmt.Sprintf("unsupported unary operation: %s%s", e.Operator, e.OperandType)
	}
}

func (e UnsupportedUnaryOperation) String() string {
	return e.Error()
}

func NewUnsupportedUnaryOperation(tok *token.Token, operator string, operandType string) *Error {
	return NewError(tok, UnsupportedUnaryOperation{Operator: operator, OperandType: operandType})
}

type UnknownOperator struct {
	Operator    string
	LeftType    string
	RightType   string
	OperandType string
}

func (e UnknownOperator) Error() string {
	switch {
	case e.LeftType != "" && e.RightType != "":
		return fmt.Sprintf("unknown operator: %s %s", e.LeftType, e.RightType)
	case e.OperandType != "":
		return fmt.Sprintf("unknown operator: %s%s", e.Operator, e.OperandType)
	default:
		return fmt.Sprintf("unknown operator: %s", e.Operator)
	}
}

func (e UnknownOperator) String() string {
	return e.Error()
}

func NewUnknownOperator(tok *token.Token, operator string, leftType string, rightType string, operandType string) *Error {
	return NewError(tok, UnknownOperator{
		Operator:    operator,
		LeftType:    leftType,
		RightType:   rightType,
		OperandType: operandType,
	})
}

type UnusableAsHashKey struct {
	Type string
}

func (e UnusableAsHashKey) Error() string {
	return fmt.Sprintf("unusable as hash key: %s", e.Type)
}

func (e UnusableAsHashKey) String() string {
	return e.Error()
}

func NewUnusableAsHashKey(tok *token.Token, typ string) *Error {
	return NewError(tok, UnusableAsHashKey{Type: typ})
}

type IndexOperatorNotSupported struct {
	Type string
}

func (e IndexOperatorNotSupported) Error() string {
	return fmt.Sprintf("index operator not supported: %s", e.Type)
}

func (e IndexOperatorNotSupported) String() string {
	return e.Error()
}

func NewIndexOperatorNotSupported(tok *token.Token, typ string) *Error {
	return NewError(tok, IndexOperatorNotSupported{Type: typ})
}

type NotCallable struct {
	Type string
}

func (e NotCallable) Error() string {
	return fmt.Sprintf("not callable: %s", e.Type)
}

func (e NotCallable) String() string {
	return e.Error()
}

func NewNotCallable(tok *token.Token, typ string) *Error {
	return NewError(tok, NotCallable{Type: typ})
}

type BuiltinArgumentType struct {
	Builtin  string
	Expected string
	Got      string
}

func (e BuiltinArgumentType) Error() string {
	return fmt.Sprintf("argument to `%s` must be %s, got %s", e.Builtin, e.Expected, e.Got)
}

func (e BuiltinArgumentType) String() string {
	return e.Error()
}

func NewBuiltinArgumentType(tok *token.Token, builtin string, expected string, got string) *Error {
	return NewError(tok, BuiltinArgumentType{Builtin: builtin, Expected: expected, Got: got})
}

type ExpectedType struct {
	Context  string
	Expected string
	Got      string
}

func (e ExpectedType) Error() string {
	return fmt.Sprintf("%s must be %s, got %s", e.Context, e.Expected, e.Got)
}

func (e ExpectedType) String() string {
	return e.Error()
}

func NewExpectedType(tok *token.Token, context string, expected string, got string) *Error {
	return NewError(tok, ExpectedType{Context: context, Expected: expected, Got: got})
}

func NewError(tok *token.Token, detail error) *Error {
	err := &Error{Detail: detail}
	if tok != nil {
		cloned := *tok
		err.Tok = &cloned
	}
	if detail != nil {
		err.Message = detail.Error()
	}
	return err
}

type UnknownNode struct {
	Node *ast.Node
}

func (e UnknownNode) Error() string {
	return fmt.Sprintf("unknown node: %T", *e.Node)
}

func (e UnknownNode) String() string { return e.Error() }

func NewUnknownNode(tok *token.Token, node *ast.Node) *Error {
	return NewError(tok, UnknownNode{Node: node})
}

type ParserErrors struct {
	Errors []parser.ParserError
}

func (e ParserErrors) Error() string {
	messages := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("parser errors:\n%s", strings.Join(messages, "\n"))
}

func (e ParserErrors) String() string { return e.Error() }

func NewParserErrors(tok *token.Token, errors []parser.ParserError) *Error {
	return NewError(tok, ParserErrors{Errors: errors})
}
