package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	for _, tc := range []struct {
		position int
		length   int
		expected string
	}{
		{0, 5, "[0:5]"},
		{100, 50, "[100:150]"},
	} {
		textSpan := NewTextSpan(tc.position, tc.length)
		require.Equal(t, tc.expected, textSpan.ToString())
	}
}

func TestBetween(t *testing.T) {
	for _, tc := range []struct {
		position         int
		length           int
		otherPosition    int
		otherLength      int
		expectedPosition int
		expectedLength   int
	}{
		{0, 5, 10, 5, 0, 15},
		{10, 5, 0, 5, 0, 15},
		{0, 5, 5, 5, 0, 10},
		{5, 5, 0, 5, 0, 10},
		{0, 5, 0, 5, 0, 5},
	} {
		ts := NewTextSpan(tc.position, tc.length)
		other := NewTextSpan(tc.otherPosition, tc.otherLength)
		between := ts.Between(other)
		require.Equal(t, tc.expectedPosition, between.Position)
		require.Equal(t, tc.expectedLength, between.Length)
	}
}
