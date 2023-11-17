package lexer

import (
	"bicep-go/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
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
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	// check if a single character "l.ch" is uniquly identifiable as a token
	for ch, tokenType := range token.UniqueSingleCharacterTokens {
		if l.ch != ch {
			continue
		}
		tok = newToken(tokenType, l.ch)
		l.readChar()
		return tok
	}

	switch l.ch {
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			tok = newToken(token.DOUBLEQUESTION, "??")
		} else {
			tok = newToken(token.QUESTION, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.EQUALS, "==")
		} else if l.peekChar() == '~' {
			l.readChar()
			tok = newToken(token.EQUALSINSENSITIVE, "=~")
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = newToken(token.ARROW, "=>")
		} else {
			tok = newToken(token.ASSIGNMENT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.NOTEQUALS, "!=")
		} else if l.peekChar() == '~' {
			l.readChar()
			tok = newToken(token.NOTEQUALSINSENSITIVE, "!~")
		} else {
			tok = newToken(token.EXCLAMATION, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.LESSTHANOREQUAL, "<=")
		} else {
			tok = newToken(token.LESSTHAN, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.GREATERTHANOREQUAL, ">=")
		} else {
			tok = newToken(token.GREATERTHAN, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = newToken(token.LOGICALOR, "||")
		} else {
			tok = newToken(token.PIPE, l.ch)
		}
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			tok = newToken(token.DOUBLECOLON, "::")
		} else {
			tok = newToken(token.COLON, l.ch)
		}

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readInteger()
			tok.Type = token.INTEGER
			return tok
		} else if l.ch == '\'' {
			tok = l.readStringToken()
		} else {
			tok = newToken(token.UNRECOGNIZED, l.ch)
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

func (l *Lexer) readInteger() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readStringToken() token.Token {
	// TODO: multiple string supports
	position := l.position
	for l.ch != 0 && l.ch != '\n' && (l.ch != '\r' || l.peekChar() != '\n') {
		l.readChar()
		if l.ch == '\'' {
			l.readChar()
			return token.Token{
				Type:    token.STRINGCOMPLETE,
				Literal: l.input[position:l.position],
			}
		}
	}

	return token.Token{
		Type:    token.STRINGLEFTPIECE,
		Literal: l.input[position:l.position],
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, literal interface{}) token.Token {
	switch v := literal.(type) {
	case []byte:
		return token.Token{Type: tokenType, Literal: string(v)}
	case byte:
		return token.Token{Type: tokenType, Literal: string(v)}
	case string:
		return token.Token{Type: tokenType, Literal: v}
	}
	return token.Token{Type: tokenType, Literal: "Unknown"}
}
