package token

import "fmt"

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	ITERATE  = ".."
	MOD      = "%"

	// Comparison
	EQUALS     = "=="
	BANG       = "!"
	NOT_EQUALS = "!="
	LESS       = "<"
	GREATER    = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	MUT      = "MUT"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	IF       = "if"
	ELSE     = "else"
	RETURN   = "return"
	FOR      = "FOR"
	IN       = "IN"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Offset  int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %s}", t.Type, t.Literal)
}

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"mut":      MUT,
	"for":      FOR,
	"in":       IN,
	"break":    BREAK,
	"continue": CONTINUE,
	"null":     NULL,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}
