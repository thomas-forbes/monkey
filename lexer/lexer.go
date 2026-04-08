package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // TODO: turn this into a rune to support UTF-8
	line         int
	column       int
	nextLine     int
	nextColumn   int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:      input,
		nextLine:   1,
		nextColumn: 1,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
		l.line = l.nextLine
		l.column = l.nextColumn
	} else {
		l.ch = l.input[l.readPosition]
		l.position = l.readPosition
		l.line = l.nextLine
		l.column = l.nextColumn
	}
	l.readPosition += 1

	if l.ch == '\n' {
		l.nextLine += 1
		l.nextColumn = 1
	} else {
		l.nextColumn += 1
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	startOffset := l.position
	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQUALS, Literal: "==", Line: startLine, Column: startColumn, Offset: startOffset}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startLine, startColumn, startOffset)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startLine, startColumn, startOffset)
	case ':':
		tok = newToken(token.COLON, l.ch, startLine, startColumn, startOffset)
	case '(':
		tok = newToken(token.LPAREN, l.ch, startLine, startColumn, startOffset)
	case ')':
		tok = newToken(token.RPAREN, l.ch, startLine, startColumn, startOffset)
	case ',':
		tok = newToken(token.COMMA, l.ch, startLine, startColumn, startOffset)
	case '+':
		tok = newToken(token.PLUS, l.ch, startLine, startColumn, startOffset)
	case '{':
		tok = newToken(token.LBRACE, l.ch, startLine, startColumn, startOffset)
	case '}':
		tok = newToken(token.RBRACE, l.ch, startLine, startColumn, startOffset)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, startLine, startColumn, startOffset)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, startLine, startColumn, startOffset)
	case '-':
		tok = newToken(token.MINUS, l.ch, startLine, startColumn, startOffset)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch, startLine, startColumn, startOffset)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startLine, startColumn, startOffset)
	case '%':
		tok = newToken(token.MOD, l.ch, startLine, startColumn, startOffset)
	case '<':
		tok = newToken(token.LESS, l.ch, startLine, startColumn, startOffset)
	case '>':
		tok = newToken(token.GREATER, l.ch, startLine, startColumn, startOffset)
	case '"':
		tok = token.Token{Type: token.STRING, Literal: l.readString(), Line: startLine, Column: startColumn, Offset: startOffset}
	case '.':
		if l.peekChar() == '.' {
			l.readChar()
			tok = token.Token{Type: token.ITERATE, Literal: "..", Line: startLine, Column: startColumn, Offset: startOffset}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn, startOffset)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQUALS, Literal: "!=", Line: startLine, Column: startColumn, Offset: startOffset}
		} else {
			tok = newToken(token.BANG, l.ch, startLine, startColumn, startOffset)
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startColumn, Offset: startOffset}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			tok.Line = startLine
			tok.Column = startColumn
			tok.Offset = startOffset
			return tok
		} else if isInt(l.ch) {
			tok.Literal = l.readInt()
			tok.Type = token.INT
			tok.Line = startLine
			tok.Column = startColumn
			tok.Offset = startOffset
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn, startOffset)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) Tokenize() []token.Token {
	tokens := []token.Token{}
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			return tokens
		}
	}
}

func (l *Lexer) readIdentifier() string {
	init_position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[init_position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readInt() string {
	init_position := l.position
	for {
		if !isInt(l.ch) {
			break
		}
		if l.ch == '.' && l.peekChar() == '.' {
			break
		}
		l.readChar()
	}
	return l.input[init_position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte, line int, column int, offset int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column, Offset: offset}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isInt(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
