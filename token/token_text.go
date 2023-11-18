package token

var tokenTypeToText = map[TokenType]string{
	TokenTypeAt:                   "@",
	TokenTypeLeftParen:            "(",
	TokenTypeRightParen:           ")",
	TokenTypeLeftBrace:            "{",
	TokenTypeRightBrace:           "}",
	TokenTypeLeftSquare:           "[",
	TokenTypeRightSquare:          "]",
	TokenTypeComma:                ",",
	TokenTypeDot:                  ".",
	TokenTypeQuestion:             "?",
	TokenTypeColon:                ":",
	TokenTypeSemicolon:            ";",
	TokenTypeAssignment:           ":",
	TokenTypePlus:                 "+",
	TokenTypeMinus:                "-",
	TokenTypeAsterisk:             "*",
	TokenTypeSlash:                "/",
	TokenTypeModulo:               "%",
	TokenTypeExclamation:          "!",
	TokenTypeLessThan:             "<",
	TokenTypeGreaterThan:          ">",
	TokenTypePipe:                 "|",
	TokenTypeLessThanOrEqual:      "<:",
	TokenTypeGreaterThanOrEqual:   ">:",
	TokenTypeEquals:               "::",
	TokenTypeNotEquals:            "!:",
	TokenTypeEqualsInsensitive:    ":~",
	TokenTypeNotEqualsInsensitive: "!~",
	TokenTypeLogicalAnd:           "&&",
	TokenTypeLogicalOr:            "||",
	TokenTypeDoubleQuestion:       "??",
	TokenTypeDoubleColon:          "::",
	TokenTypeArrow:                ":>",
	TokenTypeTrueKeyword:          "true",
	TokenTypeFalseKeyword:         "false",
	TokenTypeNullKeyword:          "null",
	TokenTypeWithKeyword:          "with",
	TokenTypeAsKeyword:            "as",
	TokenTypeIdentifier:           "<Identifier>",
	TokenTypeStringLeftPiece:      "<StringLeftPiece>",
	TokenTypeStringMiddlePiece:    "<StringMiddlePiece>",
	TokenTypeStringRightPiece:     "<StringRightPiece>",
	TokenTypeStringComplete:       "<StringComplete>",
	TokenTypeMultilineString:      "<MultilineString>",
	TokenTypeUnrecognized:         "<Unrecognized>",
	TokenTypeInteger:              "<Integer>",
	TokenTypeNewLine:              "<NewLine>",
	TokenTypeEndOfFile:            "<EndOfFile>",
}

func GetTokenText(t TokenType) string {
	if value, ok := tokenTypeToText[t]; ok {
		return value
	}
	return ""
}
