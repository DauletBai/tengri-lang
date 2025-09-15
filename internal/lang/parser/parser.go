// FILE: internal/lang/parser/parser.go
// Purpose: Pratt parser for Tengri with arrays and calls, Kazakh keywords, and optional type annotations.

package parser

import (
	"fmt"
	"strconv"

	"github.com/DauletBai/tengri-lang/internal/lang/ast"
	"github.com/DauletBai/tengri-lang/internal/lang/lexer"
	"github.com/DauletBai/tengri-lang/internal/lang/token"
)

// Precedences
const (
	_ int = iota
	LOWEST
	EQUALS      // == !=
	LESSGREATER // < > <= >=
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // - !
	CALL        // f(x)
)

var precedences = map[token.TokenType]int{
	token.Op_Equal:     EQUALS,
	token.Op_NotEqual:  EQUALS,
	token.Op_Less:      LESSGREATER,
	token.Op_LessEq:    LESSGREATER,
	token.Op_Greater:   LESSGREATER,
	token.Op_GreaterEq: LESSGREATER,
	token.Op_Plus:      SUM,
	token.Op_Minus:     SUM,
	token.Op_Multiply:  PRODUCT,
	token.Op_Divide:    PRODUCT,
	token.Sep_LParen:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Initialize token window
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.IntLiteral, p.parseIntegerLiteral)
	p.registerPrefix(token.JAN, p.parseBoolean)
	p.registerPrefix(token.JYN, p.parseBoolean)
	p.registerPrefix(token.Op_Minus, p.parsePrefixExpression)
	p.registerPrefix(token.Op_Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Sep_LParen, p.parseGroupedExpression)
	p.registerPrefix(token.Sep_LBracket, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Op_Plus, p.parseInfixExpression)
	p.registerInfix(token.Op_Minus, p.parseInfixExpression)
	p.registerInfix(token.Op_Multiply, p.parseInfixExpression)
	p.registerInfix(token.Op_Divide, p.parseInfixExpression)
	p.registerInfix(token.Op_Equal, p.parseInfixExpression)
	p.registerInfix(token.Op_NotEqual, p.parseInfixExpression)
	p.registerInfix(token.Op_Less, p.parseInfixExpression)
	p.registerInfix(token.Op_LessEq, p.parseInfixExpression)
	p.registerInfix(token.Op_Greater, p.parseInfixExpression)
	p.registerInfix(token.Op_GreaterEq, p.parseInfixExpression)
	p.registerInfix(token.Sep_LParen, p.parseCallExpression)

	return p
}

func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// Program := Statement*
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.BEKIT:
		return p.parseBekitStatement()
	case token.JASA:
		return p.parseJasaStatement()
	case token.QAITAR:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseJasaStatement() *ast.JasaStatement {
	stmt := &ast.JasaStatement{Token: p.curToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: ':' Type
	if p.peekTokenIs(token.Op_Colon) {
		p.nextToken() // ':'
		p.nextToken() // type
		stmt.Type = &ast.TypeNode{Token: p.curToken}
	}

	if !p.expectPeek(token.Op_Assign) {
		return nil
	}

	p.nextToken() // move past '='
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBekitStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.Op_Colon) {
		p.nextToken()
		p.nextToken()
		stmt.Type = &ast.TypeNode{Token: p.curToken}
	}

	if !p.expectPeek(token.Op_Assign) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.JAN)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.Sep_RParen) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.Sep_RParen)
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{Token: p.curToken}
	al.Elements = p.parseExpressionList(token.Sep_RBracket)
	return al
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.Sep_Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// --- registration helpers ---
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}