// FILE: internal/lang/lexer/lexer.go
// Purpose: UTF-8 aware lexer for Tengri language. Handles Kazakh keywords and ASCII operators.
// Notes: Unicode letters allowed in identifiers; whitespace skipping uses unicode.IsSpace.

package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/DauletBai/tengri-lang/internal/lang/token"
)

type Lexer struct {
	input        string
	position     int  // current read position (byte index) of current rune
	readPosition int  // next read position (byte index)
	ch           rune // current rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances by one UTF-8 rune.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
		return
	}
	r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
	l.ch = r
	l.position = l.readPosition
	l.readPosition += size
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func isLetter(ch rune) bool {
	if ch == '_' || ch == '\'' { // allow underscore and apostrophe
		return true
	}
	return unicode.IsLetter(ch)
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	// move past opening quote
	l.readChar()
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	s := l.input[start:l.position]
	return s
}

func newToken(tt token.TokenType, ch rune) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}

// Keywords table
var keywords = map[string]token.TokenType{
	"jasa":    token.JASA,
	"bekit":   token.BEKIT,
	"atqar'm": token.ATQARM,
	"qaıtar":  token.QAITAR,
	"eger":    token.EGER,
	"áıtpece": token.AITPECE,
	"ázirshe": token.AZIRSHE,
	"jan":     token.JAN,
	"j'n":     token.JYN,
	"kórset":  token.KORSET,
	"san":     token.SAN,
	"bólshek": token.BOLSHEK,
	"jol":     token.JOL,
	"tańba":   token.TANBA,
	"aqıqat":  token.AQIQAT,
	"j'i'm":   token.JYIM,
}

func LookupIdent(ident string) token.TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return token.Identifier
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.Sep_Comma, l.ch)
	case '(':
		tok = newToken(token.Sep_LParen, l.ch)
	case ')':
		tok = newToken(token.Sep_RParen, l.ch)
	case '{':
		tok = newToken(token.Sep_LBrace, l.ch)
	case '}':
		tok = newToken(token.Sep_RBrace, l.ch)
	case '[':
		tok = newToken(token.Sep_LBracket, l.ch)
	case ']':
		tok = newToken(token.Sep_RBracket, l.ch)
	case ':':
		tok = newToken(token.Op_Colon, l.ch)
	case '+':
		tok = newToken(token.Op_Plus, l.ch)
	case '-':
		tok = newToken(token.Op_Minus, l.ch)
	case '*':
		tok = newToken(token.Op_Multiply, l.ch)
	case '/':
		tok = newToken(token.Op_Divide, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Op_NotEqual, Literal: lit}
		} else {
			tok = newToken(token.Op_Bang, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Op_Equal, Literal: lit}
		} else {
			tok = newToken(token.Op_Assign, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Op_LessEq, Literal: lit}
		} else {
			tok = newToken(token.Op_Less, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Op_GreaterEq, Literal: lit}
		} else {
			tok = newToken(token.Op_Greater, l.ch)
		}
	case '"':
		tok.Type = token.StringLiteral
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = LookupIdent(literal)
			tok.Literal = literal
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