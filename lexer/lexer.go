package lexer

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

const (
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	COLON    TokenType = ":"
	COMMA    TokenType = ","

	STRING TokenType = "STRING"
	NUMBER TokenType = "NUMBER"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	NULL   TokenType = "NULL"
	EOF    TokenType = "EOF"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = Token{Type: LBRACE, Value: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Value: string(l.ch)}
	case '[':
		tok = Token{Type: LBRACKET, Value: string(l.ch)}
	case ']':
		tok = Token{Type: RBRACKET, Value: string(l.ch)}
	case ':':
		tok = Token{Type: COLON, Value: string(l.ch)}
	case ',':
		tok = Token{Type: COMMA, Value: string(l.ch)}
	case 't':
		if l.matchLiteral("true") {
			tok = Token{Type: TRUE, Value: "true"}
		} else {
			tok = Token{Type: EOF, Value: ""}
		}
	case 'f':
		if l.matchLiteral("false") {
			tok = Token{Type: FALSE, Value: "false"}
		} else {
			tok = Token{Type: EOF, Value: ""}
		}
	case 'n':
		if l.matchLiteral("null") {
			tok = Token{Type: NULL, Value: "null"}
		} else {
			tok = Token{Type: EOF, Value: ""}
		}
	case '"':
		tok.Type = STRING
		val, e := l.readString()
		if e != nil {
			fmt.Printf("Error reading string: %v\n", e)
			tok.Type = EOF
			tok.Value = ""
		} else {
			tok.Value = val
		}
	case 0:
		tok = Token{Type: EOF, Value: ""}
	default:
		if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Value = l.readNumber()
			return tok
		} else {
			tok = Token{Type: EOF, Value: ""}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() (string, error) {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	if l.ch == 0 {
		return "", fmt.Errorf("unterminated string")
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) matchLiteral(lit string) bool {
	end := l.position + len(lit)
	if end > len(l.input) {
		return false
	}
	if l.input[l.position:end] != lit {
		return false
	}
	// Advance past the literal
	for i := 1; i < len(lit); i++ {
		l.readChar()
	}
	return true
}
