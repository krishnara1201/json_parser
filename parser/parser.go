package parser

import (
	"fmt"
	"json-parser/ast"
	"json-parser/lexer"
	"strconv"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Parse() ast.JSONValue {
	v := p.parseValue()

	// Optional root check: ensure no trailing tokens.
	if !p.curTokenIs(lexer.EOF) && !p.peekTokenIs(lexer.EOF) {
		p.errors = append(p.errors, "trailing tokens after root JSON value")
	}
	return v
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expected next token %q, got %q", t, p.peekToken.Type))
	return false
}

func (p *Parser) parseValue() ast.JSONValue {
	switch p.curToken.Type {
	case lexer.STRING:
		return ast.JSONString(p.curToken.Value)
	case lexer.NUMBER:
		n, err := strconv.ParseFloat(p.curToken.Value, 64)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("invalid number: %s", p.curToken.Value))
			return nil
		}
		return ast.JSONNumber(n)
	case lexer.TRUE:
		return ast.JSONBoolean(true)
	case lexer.FALSE:
		return ast.JSONBoolean(false)
	case lexer.NULL:
		return ast.JSONNull{}
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token: %s", p.curToken.Value))
		return nil
	}
}

func (p *Parser) parseObject() ast.JSONObject {
	obj := make(ast.JSONObject)

	if p.peekToken.Type == lexer.RBRACE {
		p.nextToken()
		p.nextToken()
		return obj
	}

	for {
		p.nextToken()
		if !p.curTokenIs(lexer.STRING) {
			p.errors = append(p.errors, fmt.Sprintf("expected string key, got: %s", p.curToken.Value))
			return nil
		}
		key := p.curToken.Value

		if !p.expectPeek(lexer.COLON) {
			p.errors = append(p.errors, fmt.Sprintf("expected colon after key, got: %s", p.peekToken.Value))
			return nil
		}
		p.nextToken()
		value := p.parseValue()
		obj[key] = value

		if p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			continue
		}
		if p.peekToken.Type == lexer.RBRACE {
			p.nextToken()
			p.nextToken()
			break
		}
		p.errors = append(p.errors, fmt.Sprintf("expected comma or closing brace, got: %s", p.peekToken.Value))
		return nil
	}
	return obj
}

func (p *Parser) parseArray() ast.JSONArray {
	var arr ast.JSONArray

	if p.peekToken.Type == lexer.RBRACKET {
		p.nextToken()
		p.nextToken()
		return arr
	}

	for {
		p.nextToken()
		value := p.parseValue()
		arr = append(arr, value)

		if p.peekToken.Type == lexer.COMMA {
			p.nextToken()
			continue
		}
		if p.peekTokenIs(lexer.RBRACKET) {
			p.nextToken()
			p.nextToken()
			break
		}

		p.errors = append(p.errors, fmt.Sprintf("expected comma or closing bracket, got: %s", p.peekToken.Value))
		return nil
	}
	return arr
}
