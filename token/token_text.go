package token

var tokenTypeToText = map[TokenType]string{
	AT:                     "@",
	LEFT_PAREN:             "(",
	RIGHT_PAREN:            ")",
	LEFT_BRACE:             "{",
	RIGHT_BRACE:            "}",
	LEFT_SQUARE:            "[",
	RIGHT_SQUARE:           "]",
	COMMA:                  ",",
	DOT:                    ".",
	QUESTION:               "?",
	COLON:                  ":",
	SEMICOLON:              ";",
	ASSIGNMENT:             ":",
	PLUS:                   "+",
	MINUS:                  "-",
	ASTERISK:               "*",
	SLASH:                  "/",
	MODULO:                 "%",
	EXCLAMATION:            "!",
	LESS_THAN:              "<",
	GREATER_THAN:           ">",
	PIPE:                   "|",
	LESS_THAN_OR_EQUAL:     "<:",
	GREATER_THAN_OR_EQUAL:  ">:",
	EQUALS:                 "::",
	NOT_EQUALS:             "!:",
	EQUALS_INSENSITIVE:     ":~",
	NOT_EQUALS_INSENSITIVE: "!~",
	LOGICAL_AND:            "&&",
	LOGICAL_OR:             "||",
	DOUBLE_QUESTION:        "??",
	DOUBLE_COLON:           "::",
	ARROW:                  ":>",
	TRUE_KEYWORD:           "true",
	FALSE_KEYWORD:          "false",
	NULL_KEYWORD:           "null",
	WITH_KEYWORD:           "with",
	AS_KEYWORD:             "as",
	IDENTIFIER:             "<IDENTIFIER>",
	STRING_LEFT_PIECE:      "<STRING_LEFT_PIECE>",
	STRING_MIDDLE_PIECE:    "<STRING_MIDDLE_PIECE>",
	STRING_RIGHT_PIECE:     "<STRING_RIGHT_PIECE>",
	STRING_COMPLETE:        "<STRING_COMPLETE>",
	MULTILINE_STRING:       "<MULTILINE_STRING>",
	UNRECOGNIZED:           "<UNRECOGNIZED>",
	INTEGER:                "<INTEGER>",
	NEW_LINE:               "<NEW_LINE>",
	END_OF_FILE:            "<END_OF_FILE>",
}

func GetTokenText(t TokenType) string {
	if value, ok := tokenTypeToText[t]; ok {
		return value
	}
	return ""
}
