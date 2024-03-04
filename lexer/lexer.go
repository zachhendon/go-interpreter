package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	input        string
	filename     string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string, filename string) *Lexer {
	l := &Lexer{
		input:    input,
		filename: filename,
		line:     1,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
	l.column += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, l *Lexer) token.Token {
	return token.Token{
		Type:     tokenType,
		Literal:  string(l.ch),
		Filename: l.filename,
		Line:     l.line,
		Column:   l.column,
	}
}

func newTwoCharToken(tokenType token.TokenType, doubleTokenType token.TokenType, expected byte, l *Lexer) token.Token {
	col := l.column
	if l.peekChar() == expected {
		ch := l.ch
		l.readChar()
		return token.Token{
			Type:     doubleTokenType,
			Literal:  string(ch) + string(l.ch),
			Filename: l.filename,
			Line:     l.line,
			Column:   col,
		}
	} else {
		return token.Token{
			Type:     tokenType,
			Literal:  string(l.ch),
			Filename: l.filename,
			Line:     l.line,
			Column:   col,
		}
	}
}

func newNumberToken(l *Lexer) token.Token {
	var tokenType token.TokenType = token.INT
	position := l.position
	col := l.column
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
		if l.ch == '.' {
			if tokenType == token.FLOAT {
				break
			}
			tokenType = token.FLOAT
		}
	}
	literal := l.input[position:l.position]
	return token.Token{
		Type:     tokenType,
		Literal:  literal,
		Filename: l.filename,
		Line:     l.line,
		Column:   col,
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
			l.column = 0
		}
		if l.ch == '\t' {
			l.column += 4
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newTwoCharToken(token.ASSIGN, token.EQ, '=', l)
	case '+':
		tok = newToken(token.PLUS, l)
	case '-':
		tok = newToken(token.MINUS, l)
	case '!':
		tok = newTwoCharToken(token.BANG, token.NOT_EQ, '=', l)
	case '*':
		tok = newToken(token.ASTERISK, l)
	case '/':
		tok = newToken(token.SLASH, l)
	case '<':
		tok = newTwoCharToken(token.LT, token.LT_EQ, '=', l)
	case '>':
		tok = newTwoCharToken(token.GT, token.GT_EQ, '=', l)
	case ';':
		tok = newToken(token.SEMICOLON, l)
	case '(':
		tok = newToken(token.LPAREN, l)
	case ')':
		tok = newToken(token.RPAREN, l)
	case ',':
		tok = newToken(token.COMMA, l)
	case '{':
		tok = newToken(token.LBRACE, l)
	case '}':
		tok = newToken(token.RBRACE, l)
	case 0:
		tok = token.Token{
			Type:     token.EOF,
			Literal:  "",
			Filename: l.filename,
			Line:     l.line,
			Column:   l.column,
		}
	default:
		if isLetter(l.ch) {
			col := l.column
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = token.Token{
				Type:     tokType,
				Literal:  literal,
				Filename: l.filename,
				Line:     l.line,
				Column:   col,
			}
			return tok
		} else if isDigit(l.ch) {
			tok = newNumberToken(l)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l)
		}
	}

	l.readChar()
	return tok
}
