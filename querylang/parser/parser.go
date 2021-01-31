package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/HawkBrave/Emfetch/querylang/ast"
	"github.com/HawkBrave/Emfetch/querylang/lexer"
	"github.com/HawkBrave/Emfetch/querylang/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenTypeIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenTypeIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseString() string {
	return p.currentToken.Literal
}

func (p *Parser) parseValue() ast.Value {
	switch p.currentToken.Type {
	case token.LeftBrace:
		return p.parseObject()
	case token.LeftBracket:
		return p.parseArray()
	default:
		return p.parseLiteral()
	}
}

func (p *Parser) parseLiteral() ast.Literal {
	val := ast.Literal{Type: "Literal"}

	defer p.nextToken()

	switch p.currentToken.Type {
	case token.String:
		val.Value = p.parseString()
		return val
	case token.Number:
		v, _ := strconv.Atoi(p.currentToken.Literal)
		val.Value = v
		return val
	case token.True:
		val.Value = true
		return val
	case token.False:
		val.Value = false
		return val
	default:
		val.Value = "null"
		return val
	}
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: "Property"}
	propState := ast.PropertyStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch propState {
		case ast.PropertyStart:
			if p.currentTokenTypeIs(token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
				}
				prop.Key = key
				propState = ast.PropertyKey
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing property start, Expected string got %v", p.currentToken.Literal))
			}
		case ast.PropertyKey:
			if p.currentTokenTypeIs(token.Colon) {
				propState = ast.PropertyColon
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing property, Expected colon got %v", p.currentToken.Literal))
			}
		case ast.PropertyColon:
			val := p.parseValue()
			prop.Value = val
			return prop
		}
	}
	return prop
}

func (p *Parser) parseObject() ast.Value {
	obj := ast.Object{Type: "Object"}
	objState := ast.ObjStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenTypeIs(token.LeftBrace) {
				objState = ast.ObjOpen
				obj.Start = p.currentToken.Start
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing object, Expected '{' got %v", p.currentToken.Literal))
				return nil
			}
		case ast.ObjOpen:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		case ast.ObjProperty:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.Start
				return obj
			} else if p.currentTokenTypeIs(token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing property Expected RightBrace or Comma token, got %v", p.currentToken.Literal))
				return nil
			}
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
		}
	}
	obj.End = p.currentToken.Start
	return obj
}

func (p *Parser) parseArray() ast.Value {
	arr := ast.Array{Type: "Array"}
	arrState := ast.ArrayStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch arrState {
		case ast.ArrayStart:
			if p.currentTokenTypeIs(token.LeftBracket) {
				arr.Start = p.currentToken.Start
				arrState = ast.ArrayOpen
				p.nextToken()
			}
		case ast.ArrayOpen:
			if p.currentTokenTypeIs(token.RightBracket) {
				arr.End = p.currentToken.End
				p.nextToken()
				return arr
			}
			val := p.parseValue()
			arr.Children = append(arr.Children, val)
			arrState = ast.ArrayValue
			if p.peekTokenTypeIs(token.RightBracket) {
				p.nextToken()
			}
		case ast.ArrayValue:
			if p.currentTokenTypeIs(token.RightBracket) {
				arr.End = p.currentToken.End
				p.nextToken()
				return arr
			} else if p.currentTokenTypeIs(token.Comma) {
				arrState = ast.ArrayComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing array, Expected RightBracket or Comma got %v", p.currentToken.Literal))
			}
		case ast.ArrayComma:
			val := p.parseValue()
			arr.Children = append(arr.Children, val)
			arrState = ast.ArrayValue
		}
	}
	arr.End = p.currentToken.Start
	return arr
}

func (p *Parser) ParseProgram() (ast.RootNode, error) {
	var rootNode ast.RootNode
	if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}
	val := p.parseValue()
	if val == nil {
		p.parseError(fmt.Sprintf("Error parsing expected a value, got %v:", p.currentToken.Literal))
		return ast.RootNode{}, errors.New(p.Errors())
	}
	rootNode.RootValue = &val
	return rootNode, nil
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.nextToken()
	p.nextToken()

	return p
}
