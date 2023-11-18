package token

import "fmt"

type TokenType int

type Token struct {
	Type           TokenType
	Literal        string
	Line           int
	leadingTrivia  []Trivia
	trailingTrivia []Trivia
}

func New(
	tokenType TokenType,
	literal string,
	leadingTrivia []Trivia,
	trailingTrivia []Trivia,
) Token {

	return Token{
		Type:           tokenType,
		Literal:        literal,
		Line:           0,
		leadingTrivia:  leadingTrivia,
		trailingTrivia: trailingTrivia,
	}
}

func (tok *Token) ToString() string {
	return fmt.Sprintf("Type: %s, Literal: %s", GetTokenText(tok.Type), tok.Literal)
}

// https://github.com/Azure/bicep/blob/main/src/Bicep.Core/Parsing/TokenType.cs
const (
	UNRECOGNIZED TokenType = iota
	AT
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_SQUARE
	RIGHT_SQUARE
	COMMA
	DOT
	QUESTION
	COLON
	SEMICOLON
	ASSIGNMENT
	PLUS
	MINUS
	ASTERISK
	SLASH
	MODULO
	EXCLAMATION
	LESS_THAN
	GREATER_THAN
	LESS_THAN_OR_EQUAL
	GREATER_THAN_OR_EQUAL
	EQUALS
	NOT_EQUALS
	EQUALS_INSENSITIVE
	NOT_EQUALS_INSENSITIVE
	LOGICAL_AND
	LOGICAL_OR
	IDENTIFIER
	STRING_LEFT_PIECE
	STRING_MIDDLE_PIECE
	STRING_RIGHT_PIECE
	STRING_COMPLETE
	MULTILINE_STRING
	INTEGER
	TRUE_KEYWORD
	FALSE_KEYWORD
	NULL_KEYWORD
	NEW_LINE
	END_OF_FILE
	DOUBLE_QUESTION
	DOUBLE_COLON
	ARROW
	PIPE
	WITH_KEYWORD
	AS_KEYWORD
)
