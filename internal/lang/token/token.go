// FILE: internal/lang/token/token.go
// Purpose: Token types for Tengri language with Kazakh keywords and operator set.
// Notes: Keep token names stable across lexer/parser/evaluator.

package token

import "fmt"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type:%s, Literal:`%s`}", t.Type, t.Literal)
}

// Token constants.
const (
	// Special
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	SEMICOLON = "SEMICOLON"

	// Identifiers & literals
	Identifier    = "Identifier"
	IntLiteral    = "IntLiteral"
	StringLiteral = "StringLiteral"

	// Keywords (Kazakh)
	JASA   = "jasa"     // let
	BEKIT  = "bekit"    // set/update
	ATQARM = "atqar'm"  // function
	QAITAR = "qaıtar"   // return
	EGER   = "eger"     // if
	AITPECE= "áıtpece"  // else
	AZIRSHE= "ázirshe"  // while / for-now (reserved)
	JAN    = "jan"      // true
	JYN    = "j'n"      // false
	KORSET = "kórset"   // print (builtin)
	SAN    = "san"      // type: integer
	BOLSHEK= "bólshek"  // type: float (reserved)
	JOL    = "jol"      // string
	TANBA  = "tańba"    // rune/char (reserved)
	AQIQAT = "aqıqat"   // type: boolean
	JYIM   = "j'i'm"    // array

	// Operators
	Op_Assign    = "="
	Op_Colon     = ":"
	Op_Plus      = "+"
	Op_Minus     = "-"
	Op_Multiply  = "*"
	Op_Divide    = "/"
	Op_Equal     = "=="
	Op_NotEqual  = "!="
	Op_Less      = "<"
	Op_LessEq    = "<="
	Op_Greater   = ">"
	Op_GreaterEq = ">="
	Op_Bang      = "!"

	// Delimiters
	Sep_LParen   = "("
	Sep_RParen   = ")"
	Sep_LBrace   = "{"
	Sep_RBrace   = "}"
	Sep_LBracket = "["
	Sep_RBracket = "]"
	Sep_Comma    = ","
)