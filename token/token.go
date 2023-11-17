package token

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"true":  TRUEKEYWORD,
	"false": FALSEKEYWORD,
	"null":  NULLKEYWORD,
	"as":    ASKEYWORD,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
