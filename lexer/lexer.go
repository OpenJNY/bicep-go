package lexer

import (
	"bicep-go/syntax"
	"bicep-go/token"
	"bicep-go/util"
	"strings"
)

var SingleCharacterEscapes = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
	'\'': '\'',
	'$':  '$',
}

const (
	MultilineStringTerminatingQuoteCount = 3
)

type Lexer struct {
	textWindow    *TextWindow
	tokens        []*token.Token
	templateStack *util.Stack[token.TokenType]
}

func New(input string) *Lexer {
	return &Lexer{
		textWindow:    NewTextWindow(input),
		tokens:        []*token.Token{},
		templateStack: util.NewStack[token.TokenType](),
	}
}

func (l *Lexer) GetTokens() []*token.Token {
	return l.tokens
}

func (l *Lexer) Lex() {
	for !l.textWindow.IsAtEnd() {
		l.LexToken()
	}

	// make sure the last token is EOF
	if len(l.tokens) == 0 || l.tokens[len(l.tokens)-1].Type != token.TokenTypeEndOfFile {
		l.LexToken()
	}
}

func (l *Lexer) LexToken() {
	l.textWindow.Reset()
	leadingTrivia := l.scanLeadingTrivia()

	l.textWindow.Reset()
	tokenType := l.scanToken()
	tokenText := l.textWindow.GetText()

	l.textWindow.Reset()
	includeComments := syntax.GetCommentStickiness(tokenType) >= syntax.COMMENT_STICKINESS_TRAILING
	trailingTrivia := l.scanTrailingTrivia(includeComments)

	token := token.NewToken(tokenType, tokenText, leadingTrivia, trailingTrivia)
	l.tokens = append(l.tokens, token)
}

var uniqueSingleCharacterTokens = map[byte]token.TokenType{
	'(': token.TokenTypeLeftParen,
	')': token.TokenTypeRightParen,
	'[': token.TokenTypeLeftSquare,
	']': token.TokenTypeRightSquare,
	'@': token.TokenTypeAt,
	',': token.TokenTypeComma,
	'.': token.TokenTypeDot,
	';': token.TokenTypeSemicolon,
	'+': token.TokenTypePlus,
	'-': token.TokenTypeMinus,
	'%': token.TokenTypeModulo,
	'*': token.TokenTypeAsterisk,
	'/': token.TokenTypeSlash,
}

func (l *Lexer) scanToken() token.TokenType {
	if l.textWindow.IsAtEnd() {
		return token.TokenTypeEndOfFile
	}

	nextChar := l.textWindow.Peek()
	l.textWindow.Advance()

	// if the token type can be identified by a single character (e.g. +), return it
	if tokenType, ok := uniqueSingleCharacterTokens[nextChar]; ok {
		return tokenType
	}

	switch nextChar {
	case '{':
		if l.templateStack.Any() {
			l.templateStack.Push(token.TokenTypeLeftBrace)
		}
		return token.TokenTypeLeftBrace
	case '}':
		if l.templateStack.Any() {
			prevTemplateToken, _ := l.templateStack.Peek()
			if prevTemplateToken == token.TokenTypeLeftBrace {
				stringToken := l.scanStringSegment(false)
				if stringToken == token.TokenTypeRightBrace {
					l.templateStack.Pop()
				}
				return stringToken
			}
		}
		return token.TokenTypeRightBrace
	case '?':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '?' {
			l.textWindow.Advance()
			return token.TokenTypeDoubleQuestion
		}
		return token.TokenTypeQuestion
	case ':':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == ':' {
			l.textWindow.Advance()
			return token.TokenTypeDoubleColon
		}
		return token.TokenTypeColon
	case '!':
		if !l.textWindow.IsAtEnd() {
			if l.textWindow.Peek() == '=' {
				l.textWindow.Advance()
				return token.TokenTypeNotEquals
			} else if l.textWindow.Peek() == '~' {
				l.textWindow.Advance()
				return token.TokenTypeNotEqualsInsensitive
			}
		}
		return token.TokenTypeExclamation
	case '<':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '=' {
			l.textWindow.Advance()
			return token.TokenTypeLessThanOrEqual
		}
		return token.TokenTypeLessThan
	case '>':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '=' {
			l.textWindow.Advance()
			return token.TokenTypeGreaterThanOrEqual
		}
		return token.TokenTypeGreaterThan
	case '=':
		if !l.textWindow.IsAtEnd() {
			switch l.textWindow.Peek() {
			case '=':
				l.textWindow.Advance()
				return token.TokenTypeEquals
			case '~':
				l.textWindow.Advance()
				return token.TokenTypeEqualsInsensitive
			case '>':
				l.textWindow.Advance()
				return token.TokenTypeArrow
			}
		}
		return token.TokenTypeAssignment
	case '&':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '&' {
			l.textWindow.Advance()
			return token.TokenTypeLogicalAnd
		}
		return token.TokenTypeUnrecognized
	case '|':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '|' {
			l.textWindow.Advance()
			return token.TokenTypeLogicalOr
		}
		return token.TokenTypePipe
	case '\'':
		if l.textWindow.Peek() == '\'' && l.textWindow.PeekAt(1) == '\'' {
			l.textWindow.AdvanceTo(2)
			return l.scanMultilineString()
		}
		tokenType := l.scanStringSegment(true)
		if tokenType == token.TokenTypeStringLeftPiece {
			l.templateStack.Push(tokenType)
		}
		return tokenType
	default:
		if isNewLine(nextChar) {
			if l.templateStack.Any() {
				// need to re-check the newline token on next pass
				l.textWindow.Rewind()
				l.templateStack = util.NewStack[token.TokenType]()
				return token.TokenTypeStringRightPiece
			}
			l.scanNewLine()
			return token.TokenTypeNewLine
		} else if isDigit(nextChar) {
			l.scanNumber()
			return token.TokenTypeInteger
		} else if isIdentifierStart(nextChar) {
			return l.scanIdentifier()
		} else {
			return token.TokenTypeUnrecognized
		}
	}
}

func (l *Lexer) scanIdentifier() token.TokenType {
	for {
		if l.textWindow.IsAtEnd() || !isIdentifierContinuation(l.textWindow.Peek()) {
			identifier := l.textWindow.GetText()

			var keywords = map[string]token.TokenType{
				"true":  token.TokenTypeTrueKeyword,
				"false": token.TokenTypeFalseKeyword,
				"null":  token.TokenTypeNullKeyword,
				"with":  token.TokenTypeWithKeyword,
				"as":    token.TokenTypeAsKeyword,
			}
			if tokenType, ok := keywords[identifier]; ok {
				return tokenType
			}

			// identifier too long
			// if len(identifier) > common.MAX_IDENTIFIER_LENGTH {
			// }
			return token.TokenTypeIdentifier
		}
		l.textWindow.Advance()
	}
}

func (l *Lexer) scanStringSegment(isAtStartOfString bool) token.TokenType {
	for {
		if l.textWindow.IsAtEnd() {
			if isAtStartOfString {
				return token.TokenTypeStringComplete
			} else {
				return token.TokenTypeStringRightPiece
			}
		}

		nextChar := l.textWindow.Peek()
		if isNewLine(nextChar) {
			if isAtStartOfString {
				return token.TokenTypeStringComplete
			} else {
				return token.TokenTypeStringRightPiece
			}
		}

		// escapeBeginPosition := l.textWindow.GetAbsolutePosition()
		l.textWindow.Advance()

		if nextChar == '\'' {
			if isAtStartOfString {
				return token.TokenTypeStringComplete
			} else {
				return token.TokenTypeStringRightPiece
			}
		}

		if nextChar == '&' && !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '{' {
			l.textWindow.Advance()
			if isAtStartOfString {
				return token.TokenTypeStringLeftPiece
			} else {
				return token.TokenTypeStringMiddlePiece
			}
		}

		// the below section is for handling escape sequences
		if nextChar != '\\' {
			continue
		}

		// <'> + <EOF>
		if l.textWindow.IsAtEnd() {
			// UnterminatedStringEscapeSequenceAtEof
			if isAtStartOfString {
				return token.TokenTypeStringComplete
			} else {
				return token.TokenTypeStringRightPiece
			}
		}

		nextChar = l.textWindow.Peek()
		l.textWindow.Advance()

		if nextChar == 'u' {
			// unicode escape sequence = \u...

			if l.textWindow.IsAtEnd() {
				continue
			}

			nextChar = l.textWindow.Peek()
			if nextChar != '{' {
				// \u must be followed by {, but it's not
				continue
			}

			l.textWindow.Advance()
			if l.textWindow.IsAtEnd() {
				// string was prematurely terminated
				// reusing the first check in the loop body to produce the diagnostic
				continue
			}

			codePointText := scanHexNumber(l.textWindow)
			if l.textWindow.IsAtEnd() {
				continue
			}

			if len(codePointText) == 0 {
				// didn't get any hex digits
				continue
			}

			nextChar = l.textWindow.Peek()
			if nextChar != '}' {
				// hex digits myust be followed by }, but it's not
				continue
			}

			l.textWindow.Advance()
			if _, err := parseCodePoint(codePointText); err != nil {
				// faild to parse the code point
				continue
			}
		} else {
			if _, ok := SingleCharacterEscapes[nextChar]; !ok {
				// invalid escape sequence
				continue
			}
		}
	}
}

func (l *Lexer) scanMultilineString() token.TokenType {
	var successiveQuotes int = 0

	for !l.textWindow.IsAtEnd() {
		nextChar := l.textWindow.Peek()
		l.textWindow.Advance()

		switch nextChar {
		case '\'':
			successiveQuotes++
			if successiveQuotes == MultilineStringTerminatingQuoteCount {
				for l.textWindow.Peek() == '\'' {
					l.textWindow.Advance()
				}
				return token.TokenTypeMultilineString
			}
		default:
			successiveQuotes = 0
			break
		}
	}

	// unterminated multi-line string
	return token.TokenTypeMultilineString
}

func (l *Lexer) scanLeadingTrivia() []*token.Trivia {
	var trivias []*token.Trivia

	for {
		if isWhitespace(l.textWindow.Peek()) {
			trivias = append(trivias, l.scanWhitespace())
		} else if l.textWindow.Peek() == '/' && l.textWindow.PeekAt(1) == '/' {
			trivias = append(trivias, l.scanSingleLineComment())
		} else if l.textWindow.Peek() == '/' && l.textWindow.PeekAt(1) == '*' {
			trivias = append(trivias, l.scanMultiLineComment())
		} else {
			break
		}
	}

	return trivias
}

func (l *Lexer) scanTrailingTrivia(includeComments bool) []*token.Trivia {
	var trivias []*token.Trivia

	for {
		next := l.textWindow.Peek()
		if isWhitespace(next) {
			trivias = append(trivias, l.scanWhitespace())
		} else if includeComments && next == '/' {
			nextNext := l.textWindow.PeekAt(1)
			if nextNext == '/' {
				trivias = append(trivias, l.scanSingleLineComment())
			} else if nextNext == '*' {
				trivias = append(trivias, l.scanMultiLineComment())
			} else {
				break
			}
		} else {
			break
		}
	}

	return trivias
}

func (l *Lexer) scanWhitespace() *token.Trivia {
	l.textWindow.Reset()

	for !l.textWindow.IsAtEnd() {
		nextChar := l.textWindow.Peek()
		if nextChar == ' ' || nextChar == '\t' {
			l.textWindow.Advance()
			continue
		}
		break
	}

	return token.NewTrivia(token.WhitespaceTrivia, l.textWindow.GetText(), l.textWindow.GetSpan())
}

func (l *Lexer) scanSingleLineComment() *token.Trivia {
	l.textWindow.Reset()
	l.textWindow.AdvanceTo(2)

	for !l.textWindow.IsAtEnd() {
		nextChar := l.textWindow.Peek()
		if isNewLine(nextChar) {
			break
		}
		l.textWindow.Advance()
	}

	return token.NewTrivia(token.SingleLineCommentTrivia, l.textWindow.GetText(), l.textWindow.GetSpan())
}

func (l *Lexer) scanMultiLineComment() *token.Trivia {
	l.textWindow.Reset()
	l.textWindow.AdvanceTo(2)

	for {
		if l.textWindow.IsAtEnd() {
			// unterminated multi-line comment
			break
		}
		nextChar := l.textWindow.Peek()
		l.textWindow.Advance()

		if nextChar != '*' || l.textWindow.Peek() != '/' {
			continue
		}

		if l.textWindow.IsAtEnd() {
			// unterminated multi-line comment
			break
		}

		nextChar = l.textWindow.Peek()
		l.textWindow.Advance()
		if nextChar == '/' {
			break
		}
	}

	return token.NewTrivia(token.MultiLineCommentTrivia, l.textWindow.GetText(), l.textWindow.GetSpan())
}

func (l *Lexer) scanNewLine() {
	for !l.textWindow.IsAtEnd() {
		nextChar := l.textWindow.Peek()
		if !isNewLine(nextChar) {
			return
		}
		l.textWindow.Advance()
	}
}

func (l *Lexer) scanNumber() {
	for {
		if l.textWindow.IsAtEnd() {
			return
		}
		if !isDigit(l.textWindow.Peek()) {
			return
		}
		l.textWindow.Advance()
	}
}

func scanHexNumber(textWindow *TextWindow) string {
	var builder strings.Builder
	for {
		if textWindow.IsAtEnd() {
			return builder.String()
		}
		current := textWindow.Peek()
		if !isHexDigit(current) {
			return builder.String()
		}
		builder.WriteByte(current)
		textWindow.Advance()
	}
}

func parseCodePoint(codePointText string) (string, error) {
	// TODO: implement this
	return "", nil
}

func isIdentifierStart(ch byte) bool {
	return isLetter(ch) || ch == '_'
}

func isIdentifierContinuation(ch byte) bool {
	return isIdentifierStart(ch) || isDigit(ch)
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func isNewLine(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}
