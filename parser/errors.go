package parser

import (
	"fmt"
	"monkey/token"
)

type ParserError interface {
	error
	Kind() string
	Token() token.Token
}

type ExpectedTokenError struct {
	Expected token.TokenType
	Got      token.Token
}

func (e ExpectedTokenError) Error() string {
	return fmt.Sprintf(
		"expected token %s, got %s instead at %d:%d",
		e.Expected,
		e.Got.Type,
		e.Got.Line,
		e.Got.Column,
	)
}
func (e ExpectedTokenError) Kind() string {
	return "expected_token"
}
func (e ExpectedTokenError) Token() token.Token {
	return e.Got
}

type NoPrefixParseFnError struct {
	Got token.Token
}

func (e NoPrefixParseFnError) Error() string {
	return fmt.Sprintf(
		"no prefix parse function for %s found at %d:%d",
		e.Got.Type,
		e.Got.Line,
		e.Got.Column,
	)
}
func (e NoPrefixParseFnError) Kind() string {
	return "no_prefix_parse_fn"
}
func (e NoPrefixParseFnError) Token() token.Token {
	return e.Got
}

type InvalidIntegerLiteralError struct {
	Tok     token.Token
	Literal string
}

func (e InvalidIntegerLiteralError) Error() string {
	return fmt.Sprintf(
		"could not parse %q as integer at %d:%d",
		e.Literal,
		e.Tok.Line,
		e.Tok.Column,
	)
}
func (e InvalidIntegerLiteralError) Kind() string {
	return "invalid_integer_literal"
}
func (e InvalidIntegerLiteralError) Token() token.Token {
	return e.Tok
}

type InvalidAssignmentTargetError struct {
	Tok        token.Token
	TargetType string
}

func (e InvalidAssignmentTargetError) Error() string {
	return fmt.Sprintf(
		"expected left-hand side of assignment to be identifier, got %s at %d:%d",
		e.TargetType,
		e.Tok.Line,
		e.Tok.Column,
	)
}
func (e InvalidAssignmentTargetError) Kind() string {
	return "invalid_assignment_target"
}
func (e InvalidAssignmentTargetError) Token() token.Token {
	return e.Tok
}

type InvalidForBindingCountError struct {
	Tok   token.Token
	Count int
}

func (e InvalidForBindingCountError) Error() string {
	return fmt.Sprintf(
		"expected 1 or 2 bindings in for statement, got %d at %d:%d",
		e.Count,
		e.Tok.Line,
		e.Tok.Column,
	)
}
func (e InvalidForBindingCountError) Kind() string {
	return "invalid_for_binding_count"
}
func (e InvalidForBindingCountError) Token() token.Token {
	return e.Tok
}
