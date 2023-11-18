package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTextWindow(t *testing.T) {
	input := `one two three
four five six
seven eight
nine

ten
`

	textWindow := NewTextWindow(input)
	require.Equal(t, textWindow.GetText(), "")
}
