package lexer

import (
	"bicep-go/util"
	"strings"
)

const (
	InvalidCharacter = 0
)

type TextWindow struct {
	text     string
	position int
	offset   int
}

func NewTextWindow(text string) *TextWindow {
	return &TextWindow{
		text:     text,
		position: 0,
		offset:   0,
	}
}

func (t *TextWindow) GetText() string {
	return t.text[t.position : t.position+t.offset]
}

func (t *TextWindow) GetSpan() *util.TextSpan {
	return util.NewTextSpan(t.position, t.offset)
}

func (t *TextWindow) GetAbsolutePosition() int {
	return t.position + t.offset
}

func (t *TextWindow) IsAtEnd() bool {
	return t.position+t.offset >= len(t.text)
}

func (t *TextWindow) Peek() byte {
	return t.PeekAt(0)
}

func (t *TextWindow) PeekAt(numChars int) byte {
	pos := t.position + t.offset + numChars
	if pos >= len(t.text) {
		return InvalidCharacter
	}
	return t.text[pos]
}

func (t *TextWindow) Next() byte {
	nextChar := t.Peek()
	if nextChar == InvalidCharacter {
		t.Advance()
	}
	return nextChar
}

func (t *TextWindow) Advance() {
	t.AdvanceTo(1)
}

func (t *TextWindow) AdvanceTo(numChars int) {
	t.offset += numChars
}

func (t *TextWindow) Rewind() {
	t.RewindTo(1)
}

func (t *TextWindow) RewindTo(numChars int) {
	t.offset -= numChars
}

func (t *TextWindow) Reset() {
	t.position += t.offset
	t.offset = 0
}

func (t *TextWindow) GetTextBetweenLineStartAndCurrentPosition() string {
	textBeforePosition := t.text[0:t.position]
	indexOfPreviousNewLine := strings.LastIndexByte(textBeforePosition, '\n')
	if indexOfPreviousNewLine == -1 || t.position == 0 {
		return textBeforePosition
	}

	return t.text[indexOfPreviousNewLine+1 : t.position]
}
