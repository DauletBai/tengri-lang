// FILE: internal/lang/ast/ast.go
// Purpose: AST node definitions for Tengri language with arrays and calls supported.

package ast

import (
	"bytes"

	"github.com/DauletBai/tengri-lang/internal/lang/token"
)

// Node is the interface that all AST nodes must implement.
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

// Program is the root of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents an identifier in the code.
type Identifier struct {
	Token token.Token // token.Identifier
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// TypeNode represents a type annotation in the code.
type TypeNode struct {
	Token token.Token // type token (e.g., token.SAN)
}

func (tn *TypeNode) expressionNode()      {}
func (tn *TypeNode) TokenLiteral() string { return tn.Token.Literal }
func (tn *TypeNode) String() string       { return tn.Token.Literal }

// ConstStatement: bekit NAME [: Type]? = Expr
type ConstStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *TypeNode
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	if cs.Type != nil {
		out.WriteString(" : ")
		out.WriteString(cs.Type.String())
	}
	out.WriteString(" = ")
	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}
	return out.String()
}

// JasaStatement: jasa NAME [: Type]? = Expr
type JasaStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *TypeNode
	Value Expression
}

func (js *JasaStatement) statementNode()       {}
func (js *JasaStatement) TokenLiteral() string { return js.Token.Literal }
func (js *JasaStatement) String() string {
	var out bytes.Buffer
	out.WriteString(js.TokenLiteral() + " ")
	out.WriteString(js.Name.String())
	if js.Type != nil {
		out.WriteString(" : ")
		out.WriteString(js.Type.String())
	}
	out.WriteString(" = ")
	if js.Value != nil {
		out.WriteString(js.Value.String())
	}
	return out.String()
}

// ReturnStatement: qaıtar Expr
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

// ExpressionStatement wraps an Expression as a Statement.
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

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Boolean literal.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string {
	if b.Value { return "jan" }
	return "j'n"
}

// PrefixExpression: (! | -) Right
type PrefixExpression struct {
	Token    token.Token // operator token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression: Left Op Right
type InfixExpression struct {
	Token    token.Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral: '[' (Expr (',' Expr)*)? ']'
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, e := range al.Elements {
		if i > 0 { out.WriteString(", ") }
		out.WriteString(e.String())
	}
	out.WriteString("]")
	return out.String()
}

// CallExpression: Function '(' Arguments? ')'
type CallExpression struct {
	Token     token.Token // '('
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, a := range ce.Arguments {
		if i > 0 { out.WriteString(", ") }
		out.WriteString(a.String())
	}
	out.WriteString(")")
	return out.String()
}

// IfExpression represents an `eger` expression.
type IfExpression struct {
	Token       token.Token // The 'eger' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("eger")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("áıtpece ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// Parameter represents a function parameter.
type Parameter struct {
	Token token.Token // The identifier token
	Name  *Identifier
	Type  *TypeNode
}

func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
func (p *Parameter) String() string {
	var out bytes.Buffer
	out.WriteString(p.Name.String())
	out.WriteString(" : ")
	if p.Type != nil {
		out.WriteString(p.Type.String())
	}
	return out.String()
}

// BlockStatement represents a sequence of statements enclosed in braces.
type BlockStatement struct {
	Token      token.Token // The { token (not yet in language)
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}