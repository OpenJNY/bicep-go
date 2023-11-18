package lexer

import (
	"bicep-go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	input := `resource test 'Provider/ResourceType@version' = {
		name: 'test'
	}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "resource"},
		{token.IDENTIFIER, "test"},
		{token.STRING_COMPLETE, "'Provider/ResourceType@version'"},
		{token.ASSIGNMENT, "="},
		{token.LEFT_BRACE, "{"},
		{token.NEW_LINE, "\\n"},
		{token.IDENTIFIER, "name"},
		{token.COLON, ":"},
		{token.STRING_COMPLETE, "'test'"},
		{token.NEW_LINE, "\\n"},
		{token.RIGHT_BRACE, "}"},
		{token.NEW_LINE, "\\n"},
	}

	lexer := New(input)
	lexer.Lex()
	tokens := lexer.GetTokens()

	assert.Equal(t, len(tests), len(tokens))
	for i, tt := range tests {
		require.Equal(t, tokens[i].Type, tt.expectedType)
		require.Equal(t, tokens[i].Literal, tt.expectedLiteral)
	}
}
