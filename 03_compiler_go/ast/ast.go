// 03_compiler_go/ast/ast.go
package ast

import (
	"bytes"
	"tengri-lang/03_compiler_go/token"
	"strings"
)

// ИНТЕРФЕЙСЫ
type Node interface {
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

// КОРНЕВОЙ УЗЕЛ 
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// УЗЛЫ-КОМАНДЫ (Statements)
type ConstStatement struct {
	Token token.Token // Токен 'Λ'
	Name  *Identifier
	Value Expression
}
func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	return cs.TokenLiteral() + " " + cs.Name.String() + " : " + cs.Value.String()
}

type FunctionDefinition struct {
	Token      token.Token // Токен 'Π'
	Name       *Identifier
	Parameters []*Parameter
	Body       *BlockStatement
}
func (fd *FunctionDefinition) statementNode()       {}
func (fd *FunctionDefinition) TokenLiteral() string { return fd.Token.Literal }
func (fd *FunctionDefinition) String() string {
	params := []string{}
	for _, p := range fd.Parameters {
		params = append(params, p.String())
	}
	return fd.TokenLiteral() + " " + fd.Name.String() + "(" + strings.Join(params, ", ") + ") " + fd.Body.String()
}

type Parameter struct {
	Token token.Token // Токен типа
	Name  *Identifier
}
func (p *Parameter) String() string {
	return p.Token.Literal + " " + p.Name.String()
}

type BlockStatement struct {
	Token      token.Token // Токен '('
	Statements []Statement
}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(")")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // Токен '→'
	ReturnValue Expression
}
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return rs.TokenLiteral() + " " + rs.ReturnValue.String()
	}
	return ""
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// УЗЛЫ-ВЫРАЖЕНИЯ (Expressions) 
type Identifier struct {
	Token token.Token
	Value string
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

type CallExpression struct {
	Token     token.Token // Токен '('
	Function  Expression
	Arguments []Expression
}
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}
