// go_compiler/lexer/lexer.go
package lexer

import (
	"tengri-lang/03_compiler_go/token"
	"unicode"
)
type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
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
	l.readPosition += 1
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // Пропускаем шум ПЕРЕД каждым токеном

	switch l.ch {
	// Полный список Рун, Операторов и Разделителей
	case 'Π': tok = newToken(token.Runa_Func_Def, l.ch)
	case '—': tok = newToken(token.Runa_Var, l.ch)
	case 'Λ': tok = newToken(token.Runa_Const, l.ch)
	case 'Y': tok = newToken(token.Runa_If, l.ch)
	case 'Q': tok = newToken(token.Runa_True, l.ch)
	case 'I': tok = newToken(token.Runa_False, l.ch)
	case '↻': tok = newToken(token.Runa_Loop, l.ch)
	// case '→': tok = newToken(token.Runa_Return, l.ch)
	// case '⁞': tok = newToken(token.Runa_Log, l.ch)
	case '→': tok = newToken(token.ARROW, l.ch)
	case '⁞': tok = newToken(token.SEMICOLON, l.ch)
	case '□': tok = newToken(token.Runa_Type_Int, l.ch)
	case '⊡': tok = newToken(token.Runa_Type_Float, l.ch)
	case '∞': tok = newToken(token.Runa_Type_Str, l.ch)
	case '◇': tok = newToken(token.Runa_Type_Char, l.ch)
	case '≡': tok = newToken(token.Runa_Type_Collection, l.ch)
	case ':': tok = newToken(token.Op_Assign, l.ch)
	case '+': tok = newToken(token.Op_Plus, l.ch)
	case '-': tok = newToken(token.Op_Minus, l.ch)
	case '*': tok = newToken(token.Op_Multiply, l.ch)
	case '/': tok = newToken(token.Op_Divide, l.ch)
	case '>': tok = newToken(token.Op_Greater, l.ch)
	case ',': tok = newToken(token.Sep_Comma, l.ch)
	case '(': tok = newToken(token.Sep_LParen, l.ch)
	case ')': tok = newToken(token.Sep_RParen, l.ch)
	case '"': // Обработка строк
		tok.Type = token.StringLiteral
		tok.Literal = l.readString()
	// Конец файла
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	// Если символ не опознан 
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.Identifier
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Type = token.IntLiteral
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// readString читает все символы до закрывающей кавычки
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		// Останавливаемся, если встретили кавычку или конец файла
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) skipWhitespace() {
	for {
		if unicode.IsSpace(l.ch) {
			l.readChar()
		} else if l.ch == '/' && l.peekChar() == '/' { // <-- ИСПРАВЛЕНИЕ: Правильная обработка комментариев
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else {
			break
		}
	}
}

// остальные вспомогательные функции 
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || unicode.IsLetter(ch)
}

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}