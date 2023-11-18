package syntax

import (
	"bicep-go/token"
)

func HasFreeFromText(t token.TokenType) bool {
	freeFromTextTypes := [...]token.TokenType{
		token.NEW_LINE,
		token.IDENTIFIER,
		token.INTEGER,
		token.STRING_LEFT_PIECE,
		token.STRING_MIDDLE_PIECE,
		token.STRING_RIGHT_PIECE,
		token.STRING_COMPLETE,
		token.MULTILINE_STRING,
		token.UNRECOGNIZED,
	}

	for _, freeFromTextType := range freeFromTextTypes {
		if t == freeFromTextType {
			return true
		}
	}

	return false
}

func GetText(t token.TokenType) string {
	tokentTypeToText := map[token.TokenType]string{
		token.AT:                     "@",
		token.LEFT_PAREN:             "(",
		token.RIGHT_PAREN:            ")",
		token.LEFT_BRACE:             "{",
		token.RIGHT_BRACE:            "}",
		token.LEFT_SQUARE:            "[",
		token.RIGHT_SQUARE:           "]",
		token.COMMA:                  ",",
		token.DOT:                    ".",
		token.QUESTION:               "?",
		token.COLON:                  ":",
		token.SEMICOLON:              ";",
		token.ASSIGNMENT:             ":",
		token.PLUS:                   "+",
		token.MINUS:                  "-",
		token.ASTERISK:               "*",
		token.SLASH:                  "/",
		token.MODULO:                 "%",
		token.EXCLAMATION:            "!",
		token.LESS_THAN:              "<",
		token.GREATER_THAN:           ">",
		token.PIPE:                   "|",
		token.LESS_THAN_OR_EQUAL:     "<:",
		token.GREATER_THAN_OR_EQUAL:  ">:",
		token.EQUALS:                 "::",
		token.NOT_EQUALS:             "!:",
		token.EQUALS_INSENSITIVE:     ":~",
		token.NOT_EQUALS_INSENSITIVE: "!~",
		token.LOGICAL_AND:            "&&",
		token.LOGICAL_OR:             "||",
		token.DOUBLE_QUESTION:        "??",
		token.DOUBLE_COLON:           "::",
		token.ARROW:                  ":>",
		token.TRUE_KEYWORD:           "true",
		token.FALSE_KEYWORD:          "false",
		token.NULL_KEYWORD:           "null",
		token.WITH_KEYWORD:           "with",
		token.AS_KEYWORD:             "as",
	}

	if val, ok := tokentTypeToText[t]; ok {
		return val
	}
	return ""
}

func GetCommentStickiness(t token.TokenType) CommentStickiness {
	if t == token.NEW_LINE ||
		t == token.AT ||
		t == token.MINUS ||
		t == token.END_OF_FILE ||
		t == token.LEFT_PAREN ||
		t == token.LEFT_SQUARE ||
		t == token.LEFT_BRACE ||
		t == token.STRING_LEFT_PIECE {
		return COMMENT_STICKINESS_LEADING
	}

	if t == token.RIGHT_PAREN ||
		t == token.RIGHT_SQUARE ||
		t == token.RIGHT_BRACE ||
		t == token.STRING_RIGHT_PIECE {
		return COMMENT_STICKINESS_TRAILING
	}

	if t == token.EXCLAMATION ||
		t == token.FALSE_KEYWORD ||
		t == token.TRUE_KEYWORD ||
		t == token.NULL_KEYWORD ||
		t == token.STRING_COMPLETE ||
		t == token.INTEGER ||
		t == token.IDENTIFIER {
		return COMMENT_STICKINESS_BIDRECTIONAL
	}

	return COMMENT_STICKINESS_NONE
}
