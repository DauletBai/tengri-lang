// FILE: internal/lang/lexer/lexer.go

package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/DauletBai/tenge/internal/lang/token"
)

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
	"aqsha":   token.AQSHA, 
	"jol":     token.JOL,
	"tańba":   token.TANBA,
	"aqıqat":  token.AQIQAT,
	"j'i'm":   token.JYIM,
}

func LookupIdent(ident string) token.TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return token.IDENT
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}
}

// MODIFIED: readNumber now handles floating-point numbers.
func (l *Lexer) readNumber() (tokType token.TokenType, lit string) {
	startPosition := l.position
	tokType = token.SAN_LIT // Assume it's an integer by default

	for unicode.IsDigit(l.ch) {
		l.readChar()
	}

	// If we find a dot, it's a decimal number.
	if l.ch == '.' {
		tokType = token.AQSHA_LIT
		l.readChar() // Consume the dot
		for unicode.IsDigit(l.ch) {
			l.readChar()
		}
	}

	return tokType, l.input[startPosition:l.position]
}


func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.TANBA, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.JYIM, Literal: literal}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '*':
		tok = newToken(token.MULTIPLY, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '>':
		tok = newToken(token.GREATER, l.ch)
	case '"':
		tok.Type = token.JOL
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(l.ch) {
			// MODIFIED: Call the new readNumber function
			tok.Type, tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '\'' || unicode.IsLetter(ch)
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}