package syntax

import (
	"bicep-go/token"
)

func IsFreeFromType(t token.TokenType) bool {
	for _, freeFormType := range [...]token.TokenType{
		token.TokenTypeNewLine,
		token.TokenTypeIdentifier,
		token.TokenTypeInteger,
		token.TokenTypeStringLeftPiece,
		token.TokenTypeStringMiddlePiece,
		token.TokenTypeStringRightPiece,
		token.TokenTypeStringComplete,
		token.TokenTypeMultilineString,
		token.TokenTypeUnrecognized,
	} {
		if t == freeFormType {
			return true
		}
	}

	return false
}

func GetCommentStickiness(t token.TokenType) CommentStickiness {
	if t == token.TokenTypeNewLine ||
		t == token.TokenTypeAt ||
		t == token.TokenTypeMinus ||
		t == token.TokenTypeEndOfFile ||
		t == token.TokenTypeLeftParen ||
		t == token.TokenTypeLeftSquare ||
		t == token.TokenTypeLeftBrace ||
		t == token.TokenTypeStringLeftPiece {
		return COMMENT_STICKINESS_LEADING
	}

	if t == token.TokenTypeRightParen ||
		t == token.TokenTypeRightSquare ||
		t == token.TokenTypeRightBrace ||
		t == token.TokenTypeStringRightPiece {
		return COMMENT_STICKINESS_TRAILING
	}

	if t == token.TokenTypeExclamation ||
		t == token.TokenTypeFalseKeyword ||
		t == token.TokenTypeTrueKeyword ||
		t == token.TokenTypeNullKeyword ||
		t == token.TokenTypeStringComplete ||
		t == token.TokenTypeInteger ||
		t == token.TokenTypeIdentifier {
		return COMMENT_STICKINESS_BIDRECTIONAL
	}

	return COMMENT_STICKINESS_NONE
}
