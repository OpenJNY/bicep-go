package lexer

import (
	"bicep-go/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `resource test 'Provider/ResourceType@version' = {
		name: 'test'
	}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "resource"},
		{token.IDENTIFIER, "test"},
		{token.STRINGCOMPLETE, "'Provider/ResourceType@version'"},
		{token.ASSIGNMENT, "="},
		{token.LEFTBRACE, "{"},
		{token.IDENTIFIER, "name"},
		{token.COLON, ":"},
		{token.STRINGCOMPLETE, "'test'"},
		{token.RIGHTBRACE, "}"},
	}

	lexer := New(input)
	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
