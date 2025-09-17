// FILE: internal/lang/token/token.go

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

// All token types are now based on the tenge language keywords.
const (
	// Special Tokens
	ILLEGAL = "ILLEGAL" // Represents a token we don't know
	EOF     = "EOF"     // End of File

	// Identifiers & Literals
	IDENT     = "IDENT"      // a, myVar, etc.
	SAN_LIT   = "SAN_LIT"    // 123
	AQSHA_LIT = "AQSHA_LIT"  // 12.34
	JOL_LIT   = "JOL_LIT"    // "hello"

	// Keywords
	JASA    = "jasa"
	BEKIT   = "bekit"
	ATQARM  = "atqar'm"
	QAITAR  = "qaıtar"
	EGER    = "eger"
	AITPECE = "áıtpece"
	AZIRSHE = "ázirshe"
	JAN     = "jan"
	JYN     = "j'n"
	KORSET  = "kórset"

	// Types
	SAN    = "san"
	AQSHA  = "aqsha"
	JOL    = "jol"
	TANBA  = "tańba"
	AQIQAT = "aqıqat"
	JYIM   = "j'i'm"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	MULTIPLY  = "*"
	DIVIDE    = "/"
	EQUAL     = "=="
	GREATER   = ">"

	// Delimiters
	COMMA     = ","
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACKET  = "["
	RBRACKET  = "]"
	ARROW     = "->"
)