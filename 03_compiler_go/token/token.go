// 03_compiler_go/token/token.go
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Идентификаторы и Литералы
	Identifier = "Identifier"
	IntLiteral = "IntLiteral"
	StringLiteral = "StringLiteral"

	// Руны-Ключевые слова
	Runa_Func_Def = "Π"
	Runa_Var      = "—"
	Runa_Const    = "Λ"
	Runa_If       = "Y"
	Runa_True     = "Q"
	Runa_False    = "I"
	Runa_Loop     = "↻"
	Runa_Return   = "→"
	Runa_Log      = "⁞"

	// Руны-Типы
	Runa_Type_Int        = "□"
	Runa_Type_Float      = "⊡"
	Runa_Type_Str        = "∞"
	Runa_Type_Char       = "◇"
	Runa_Type_Collection = "≡"

	// Операторы
	Op_Assign    = ":"
	Op_Plus      = "+"
	Op_Minus     = "-"
	Op_Multiply  = "*"
	Op_Divide    = "/"
	Op_Equal     = "=="
	Op_Access_At = "@"
	Op_Get_Size  = "#"
	Op_In        = "∈"
	Op_Push      = "←"
	Op_Greater 	 = ">"

	// Разделители
	Sep_LParen   = "("
	Sep_RParen   = ")"
	Sep_LBracket = "["
	Sep_RBracket = "]"
	Sep_Comma    = ","
)
