// TODO: add Unicode & UTF-8 support
package lexer

import (
	"mky/internal/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (adter current char)
	ch           byte // current char under examination
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readToken(ifn func(byte) bool) string {
	position := l.position
	for ifn(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) makeTwoCharToken(tt token.TokenType) token.Token {
	var tok token.Token

	if l.peekChar() == '=' {
		ch := l.ch
		l.readChar()
		tok.Literal = string(ch) + string(l.ch)
		if tt == token.ASSIGN {
			tok.Type = token.EQ
		} else if tt == token.BANG {
			tok.Type = token.NOT_EQ
		} else {
			tok.Type = token.ILLEGAL
		}
	} else {
		tok = newToken(tt, l.ch)
	}
	return tok
}

func (l *Lexer) skipWhitespace() {
	for isSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = l.makeTwoCharToken(token.ASSIGN)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = l.makeTwoCharToken(token.BANG)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readToken(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readToken(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isSpace(ch byte) bool {
	for _, sp := range []byte{' ', '\f', '\n', '\r', '\t', '\v'} {
		if ch == sp {
			return true
		}
	}
	return false
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}
