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
	tokens        []token.Token
	templateStack *util.Stack[token.TokenType]
}

func New(input string) *Lexer {
	return &Lexer{
		textWindow:    NewTextWindow(input),
		tokens:        []token.Token{},
		templateStack: util.NewStack[token.TokenType](),
	}
}

func (l *Lexer) GetTokens() []token.Token {
	return l.tokens
}

func (l *Lexer) Lex() {
	for !l.textWindow.IsAtEnd() {
		l.LexToken()
	}

	// make sure the last token is EOF
	if len(l.tokens) == 0 || l.tokens[len(l.tokens)-1].Type != token.END_OF_FILE {
		l.LexToken()
	}
}

func (l *Lexer) LexToken() token.Token {
	l.textWindow.Reset()
	leadingTrivia := l.scanLeadingTrivia()

	l.textWindow.Reset()
	tokenType := l.scanToken()
	tokenText := l.textWindow.GetText()

	l.textWindow.Reset()
	includeComments := syntax.GetCommentStickiness(tokenType) >= syntax.COMMENT_STICKINESS_TRAILING
	trailingTrivia := l.scanTrailingTrivia(includeComments)

	token := token.New(tokenType, tokenText, leadingTrivia, trailingTrivia)
	l.tokens = append(l.tokens, token)
	return token
}

var uniqueSingleCharacterTokens = map[byte]token.TokenType{
	'(': token.LEFT_PAREN,
	')': token.RIGHT_PAREN,
	'[': token.LEFT_SQUARE,
	']': token.RIGHT_SQUARE,
	'@': token.AT,
	',': token.COMMA,
	'.': token.DOT,
	';': token.SEMICOLON,
	'+': token.PLUS,
	'-': token.MINUS,
	'%': token.MODULO,
	'*': token.ASTERISK,
	'/': token.SLASH,
}

func (l *Lexer) scanToken() token.TokenType {
	if l.textWindow.IsAtEnd() {
		return token.END_OF_FILE
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
			l.templateStack.Push(token.LEFT_BRACE)
		}
		return token.LEFT_BRACE
	case '}':
		if l.templateStack.Any() {
			prevTemplateToken, _ := l.templateStack.Peek()
			if prevTemplateToken == token.LEFT_BRACE {
				stringToken := l.scanStringSegment(false)
				if stringToken == token.STRING_RIGHT_PIECE {
					l.templateStack.Pop()
				}
				return stringToken
			}
		}
		return token.RIGHT_BRACE
	case '?':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '?' {
			l.textWindow.Advance()
			return token.DOUBLE_QUESTION
		}
		return token.QUESTION
	case ':':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == ':' {
			l.textWindow.Advance()
			return token.DOUBLE_COLON
		}
		return token.COLON
	case '!':
		if !l.textWindow.IsAtEnd() {
			if l.textWindow.Peek() == '=' {
				l.textWindow.Advance()
				return token.NOT_EQUALS
			} else if l.textWindow.Peek() == '~' {
				l.textWindow.Advance()
				return token.NOT_EQUALS_INSENSITIVE
			}
		}
		return token.EXCLAMATION
	case '<':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '=' {
			l.textWindow.Advance()
			return token.LESS_THAN_OR_EQUAL
		}
		return token.LESS_THAN
	case '>':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '=' {
			l.textWindow.Advance()
			return token.GREATER_THAN_OR_EQUAL
		}
		return token.GREATER_THAN
	case '=':
		if !l.textWindow.IsAtEnd() {
			switch l.textWindow.Peek() {
			case '=':
				l.textWindow.Advance()
				return token.EQUALS
			case '~':
				l.textWindow.Advance()
				return token.EQUALS_INSENSITIVE
			case '>':
				l.textWindow.Advance()
				return token.ARROW
			}
		}
		return token.ASSIGNMENT
	case '&':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '&' {
			l.textWindow.Advance()
			return token.LOGICAL_AND
		}
		return token.UNRECOGNIZED
	case '|':
		if !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '|' {
			l.textWindow.Advance()
			return token.LOGICAL_OR
		}
		return token.PIPE
	case '\'':
		if l.textWindow.Peek() == '\'' && l.textWindow.PeekAt(1) == '\'' {
			l.textWindow.AdvanceTo(2)
			return l.scanMultilineString()
		}
		tok := l.scanStringSegment(true)
		if tok == token.STRING_LEFT_PIECE {
			l.templateStack.Push(tok)
		}
		return tok
	default:
		if isNewLine(nextChar) {
			if l.templateStack.Any() {
				// need to re-check the newline token on next pass
				l.textWindow.Rewind()
				l.templateStack = util.NewStack[token.TokenType]()
				return token.STRING_RIGHT_PIECE
			}
			l.scanNewLine()
			return token.NEW_LINE
		} else if isDigit(nextChar) {
			l.scanNumber()
			return token.INTEGER
		} else if isIdentifierStart(nextChar) {
			return l.scanIdentifier()
		} else {
			return token.UNRECOGNIZED
		}
	}
}

func (l *Lexer) scanIdentifier() token.TokenType {
	for {
		if l.textWindow.IsAtEnd() || !isIdentifierContinuation(l.textWindow.Peek()) {
			identifier := l.textWindow.GetText()

			var keywords = map[string]token.TokenType{
				"true":  token.TRUE_KEYWORD,
				"false": token.FALSE_KEYWORD,
				"null":  token.NULL_KEYWORD,
				"with":  token.WITH_KEYWORD,
				"as":    token.AS_KEYWORD,
			}
			if tokenType, ok := keywords[identifier]; ok {
				return tokenType
			}

			// identifier too long
			// if len(identifier) > common.MAX_IDENTIFIER_LENGTH {
			// }
			return token.IDENTIFIER
		}
		l.textWindow.Advance()
	}
}

func (l *Lexer) scanStringSegment(isAtStartOfString bool) token.TokenType {
	for {
		if l.textWindow.IsAtEnd() {
			if isAtStartOfString {
				return token.STRING_COMPLETE
			} else {
				return token.STRING_RIGHT_PIECE
			}
		}

		nextChar := l.textWindow.Peek()
		if isNewLine(nextChar) {
			if isAtStartOfString {
				return token.STRING_COMPLETE
			} else {
				return token.STRING_RIGHT_PIECE
			}
		}

		// escapeBeginPosition := l.textWindow.GetAbsolutePosition()
		l.textWindow.Advance()

		if nextChar == '\'' {
			if isAtStartOfString {
				return token.STRING_COMPLETE
			} else {
				return token.STRING_RIGHT_PIECE
			}
		}

		if nextChar == '&' && !l.textWindow.IsAtEnd() && l.textWindow.Peek() == '{' {
			l.textWindow.Advance()
			if isAtStartOfString {
				return token.STRING_LEFT_PIECE
			} else {
				return token.STRING_MIDDLE_PIECE
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
				return token.STRING_COMPLETE
			} else {
				return token.STRING_RIGHT_PIECE
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
				return token.MULTILINE_STRING
			}
		default:
			successiveQuotes = 0
			break
		}
	}

	// unterminated multi-line string
	return token.MULTILINE_STRING
}

func (l *Lexer) scanLeadingTrivia() []token.Trivia {
	var trivias []token.Trivia

	for {
		if isWhitespace(l.textWindow.Peek()) {
			trivias = append(trivias, *l.scanWhitespace())
		} else if l.textWindow.Peek() == '/' && l.textWindow.PeekAt(1) == '/' {
			trivias = append(trivias, *l.scanSingleLineComment())
		} else if l.textWindow.Peek() == '/' && l.textWindow.PeekAt(1) == '*' {
			trivias = append(trivias, *l.scanMultiLineComment())
		} else {
			break
		}
	}

	return trivias
}

func (l *Lexer) scanTrailingTrivia(includeComments bool) []token.Trivia {
	var trivias []token.Trivia

	for {
		next := l.textWindow.Peek()
		if isWhitespace(next) {
			trivias = append(trivias, *l.scanWhitespace())
		} else if includeComments && next == '/' {
			nextNext := l.textWindow.PeekAt(1)
			if nextNext == '/' {
				trivias = append(trivias, *l.scanSingleLineComment())
			} else if nextNext == '*' {
				trivias = append(trivias, *l.scanMultiLineComment())
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
