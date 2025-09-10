// 03_compiler_go/parser/parser.go
package parser

import (
	"fmt"
	"strconv"
	"github.com/DauletBai/tengri-lang/internal/lang/ast"
	"github.com/DauletBai/tengri-lang/internal/lang/lexer"
	"github.com/DauletBai/tengri-lang/internal/lang/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	CALL
)

var precedences = map[token.TokenType]int{
	token.Op_Plus:     SUM,
	token.Op_Minus:    SUM,
	token.Op_Multiply: PRODUCT,
	token.Op_Divide:   PRODUCT,
	token.Op_Greater:  LESSGREATER,
	token.Sep_LParen:  CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	errors         []string
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.IntLiteral, p.parseIntegerLiteral)
	p.registerPrefix(token.Sep_LParen, p.parseGroupedExpression)
	p.registerPrefix(token.StringLiteral, p.parseStringLiteral)

	p.registerPrefix(token.Runa_If, p.parseIfExpression)
	p.registerPrefix(token.Runa_True, p.parseBoolean)
	p.registerPrefix(token.Runa_False, p.parseBoolean)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Op_Plus, p.parseInfixExpression)
	p.registerInfix(token.Op_Minus, p.parseInfixExpression)
	p.registerInfix(token.Op_Multiply, p.parseInfixExpression)
	p.registerInfix(token.Op_Divide, p.parseInfixExpression)
	p.registerInfix(token.Op_Greater, p.parseInfixExpression)
	p.registerInfix(token.Sep_LParen, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

// --- ИСПРАВЛЕНО: Добавлена недостающая функция ---
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.Sep_LParen) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	
	// Здесь можно будет добавить обработку `AITPESE` (else) в будущем

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.Runa_True)}
}

func (p *Parser) Errors() []string { return p.errors }
func (p *Parser) nextToken()       { p.curToken = p.peekToken; p.peekToken = p.l.NextToken() }

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{Statements: []ast.Statement{}}
    for !p.curTokenIs(token.EOF) {
        if p.curTokenIs(token.SEMICOLON) { 
            p.nextToken()
            continue
        }
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Runa_Const:
		return p.parseConstStatement()
	case token.Runa_Func_Def:
		return p.parseFunctionDefinition()
	case token.Runa_Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}
	if !p.expectPeek(token.Identifier) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
    if p.peekTokenIs(token.SEMICOLON) {
        p.nextToken()
    }
    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
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
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("не удалось преобразовать %q в число", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.Sep_RParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseFunctionDefinition() ast.Statement {
	fd := &ast.FunctionDefinition{Token: p.curToken}

	// имя функции
	if !p.expectPeek(token.Identifier) { return nil }
	fd.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// список параметров: '(' ... ')'
	if !p.expectPeek(token.Sep_LParen) { return nil }
	fd.Parameters = p.parseFunctionParameters()

	// --- короткая форма "→ expr" (неявный return) ---
	// поддержим обе ситуации: ARROW в peek или уже в cur
	if p.peekTokenIs(token.ARROW) {
		p.nextToken() // cur = ARROW
	}
	if p.curTokenIs(token.ARROW) {
		p.nextToken() // перейти к первому токену выражения
		expr := p.parseExpression(LOWEST)

		ret := &ast.ReturnStatement{
			Token:       token.Token{Type: token.ARROW, Literal: "→"},
			ReturnValue: expr,
		}
		fd.Body = &ast.BlockStatement{
			Token:      token.Token{Type: token.Sep_LParen, Literal: "("},
			Statements: []ast.Statement{ret},
		}
		return fd
	}
	// --- обычное тело "( ... )" ---
	if !p.expectPeek(token.Sep_LParen) { return nil }
	fd.Body = p.parseBlockStatement()
	return fd
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
    params := []*ast.Parameter{}

    if p.peekTokenIs(token.Sep_RParen) {
        p.nextToken()
        return params
    }

    p.nextToken()

    if p.curTokenIs(token.Sep_RParen) {
        return params
    }

    param := &ast.Parameter{Token: p.curToken}
    if !isTypeToken(param.Token.Type) {
        p.errors = append(p.errors, fmt.Sprintf("ожидался тип параметра, но получен %s", p.curToken.Type))
        return nil
    }

    if !p.expectPeek(token.Identifier) {
        return nil
    }
    param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    params = append(params, param)

    for p.peekTokenIs(token.Sep_Comma) {
        p.nextToken()
        p.nextToken()
        
        param := &ast.Parameter{Token: p.curToken}
        if !isTypeToken(param.Token.Type) {
            p.errors = append(p.errors, fmt.Sprintf("ожидался тип параметра, но получен %s", p.curToken.Type))
            return nil
        }

        if !p.expectPeek(token.Identifier) {
            return nil
        }
        param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        params = append(params, param)
    }

    if !p.expectPeek(token.Sep_RParen) {
        return nil
    }

    return params
}


func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: p.curToken}

    p.nextToken() 
    for !p.curTokenIs(token.Sep_RParen) && !p.curTokenIs(token.EOF) {
        if p.curTokenIs(token.SEMICOLON) { 
            p.nextToken()
            continue
        }
        stmt := p.parseStatement()
        if stmt != nil {
            block.Statements = append(block.Statements, stmt)
        }
        p.nextToken()
    }
    return block
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.Sep_RParen)
	return exp
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

func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("ошибка: ожидался токен %s, но получен %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("не найдена функция для разбора токена '%s'", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

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

func isTypeToken(t token.TokenType) bool {
	return t == token.Runa_Type_Int
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}