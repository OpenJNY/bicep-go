package token

import "fmt"

type Token struct {
	Type           TokenType
	Literal        string
	Line           int
	leadingTrivia  []*Trivia
	trailingTrivia []*Trivia
}

func NewToken(
	tokenType TokenType,
	literal string,
	leadingTrivia []*Trivia,
	trailingTrivia []*Trivia,
) *Token {

	return &Token{
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

type TokenType int

// https://github.com/Azure/bicep/blob/main/src/Bicep.Core/Parsing/TokenType.cs
const (
	TokenTypeAt TokenType = iota
	TokenTypeUnrecognized
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeLeftParen
	TokenTypeRightParen
	TokenTypeLeftSquare
	TokenTypeRightSquare
	TokenTypeComma
	TokenTypeDot
	TokenTypeQuestion
	TokenTypeColon
	TokenTypeSemicolon
	TokenTypeAssignment
	TokenTypePlus
	TokenTypeMinus
	TokenTypeAsterisk
	TokenTypeSlash
	TokenTypeModulo
	TokenTypeExclamation
	TokenTypeLessThan
	TokenTypeGreaterThan
	TokenTypeLessThanOrEqual
	TokenTypeGreaterThanOrEqual
	TokenTypeEquals
	TokenTypeNotEquals
	TokenTypeEqualsInsensitive
	TokenTypeNotEqualsInsensitive
	TokenTypeLogicalAnd
	TokenTypeLogicalOr
	TokenTypeIdentifier
	TokenTypeStringLeftPiece
	TokenTypeStringMiddlePiece
	TokenTypeStringRightPiece
	TokenTypeStringComplete
	TokenTypeMultilineString
	TokenTypeInteger
	TokenTypeTrueKeyword
	TokenTypeFalseKeyword
	TokenTypeNullKeyword
	TokenTypeNewLine
	TokenTypeEndOfFile
	TokenTypeDoubleQuestion
	TokenTypeDoubleColon
	TokenTypeArrow
	TokenTypePipe
	TokenTypeWithKeyword
	TokenTypeAsKeyword
)
