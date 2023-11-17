package token

type TokenType string

// https://github.com/Azure/bicep/blob/main/src/Bicep.Core/Parsing/TokenType.cs
const (
	UNRECOGNIZED = "<UNRECOGNIZED>"
	EOF          = "<EOF>"
	NEWLINE      = "<NEWLINE>"

	// single-character tokens
	AT          = "@"
	LEFTPAREN   = "("
	RIGHTPAREN  = ")"
	LEFTBRACE   = "{"
	RIGHTBRACE  = "}"
	LEFTSQUARE  = "["
	RIGHTSQUARE = "]"
	COMMA       = ","
	DOT         = "."
	QUESTION    = "?"
	COLON       = ":"
	SEMICOLON   = ";"
	ASSIGNMENT  = "="
	PLUS        = "+"
	MINUS       = "-"
	ASTERISK    = "*"
	SLASH       = "/"
	MODULO      = "%"
	EXCLAMATION = "!"
	LESSTHAN    = "<"
	GREATERTHAN = ">"
	PIPE        = "|"

	// two characters tokens
	LESSTHANOREQUAL      = "<="
	GREATERTHANOREQUAL   = ">="
	EQUALS               = "=="
	NOTEQUALS            = "!="
	EQUALSINSENSITIVE    = "=~"
	NOTEQUALSINSENSITIVE = "!~"
	LOGICALAND           = "&&"
	LOGICALOR            = "||"
	DOUBLEQUESTION       = "??"
	DOUBLECOLON          = "::"
	ARROW                = "=>"

	IDENTIFIER = "<IDENTIFIER>"

	STRINGLEFTPIECE   = "<STRINGLEFTPIECE>"
	STRINGMIDDLEPIECE = "<STRINGMIDDLEPIECE>"
	STRINGRIGHTPIECE  = "<STRINGRIGHTPIECE>"
	STRINGCOMPLETE    = "<STRINGCOMPLETE>"
	MULTILINESTRING   = "<MULTILINESTRING>"
	INTEGER           = "<INTEGER>"

	// keyword
	TRUEKEYWORD  = "true"
	FALSEKEYWORD = "false"
	NULLKEYWORD  = "null"
	WITHKEYWORD  = "with"
	ASKEYWORD    = "as"
)

var UniqueSingleCharacterTokens = map[byte]TokenType{
	'@': AT,
	'(': LEFTPAREN,
	')': RIGHTPAREN,
	'{': LEFTBRACE,
	'}': RIGHTBRACE,
	'[': LEFTSQUARE,
	']': RIGHTSQUARE,
	',': COMMA,
	'.': DOT,
	// '?': QUESTION,
	// ':': COLON,
	';': SEMICOLON,
	// '=': ASSIGNMENT,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'%': MODULO,
	// '!': EXCLAMATION,
	// '<': LESSTHAN,
	// '>': GREATERTHAN,
	// '|': PIPE,
}
