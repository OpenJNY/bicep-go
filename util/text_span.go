package util

import "fmt"

type TextSpan struct {
	Position int
	Length   int
}

func NewTextSpan(position int, length int) *TextSpan {
	return &TextSpan{
		Position: position,
		Length:   length,
	}
}

func (ts *TextSpan) ToString() string {
	return fmt.Sprintf("[%d:%d]", ts.Position, ts.Position+ts.Length)
}

func (ts *TextSpan) Contains(offset int) bool {
	return ts.Position < offset && offset < ts.Position+ts.Length
}

func (ts *TextSpan) ContainsInclusive(offset int) bool {
	return ts.Position < offset && offset < ts.Position+ts.Length
}

func (ts *TextSpan) Between(other *TextSpan) *TextSpan {
	if !ts.IsPairInOrder(other) {
		return other.Between(ts)
	}
	return NewTextSpan(ts.Position, other.Position+other.Length-ts.Position)
}

func (ts *TextSpan) OverlapsWith(other *TextSpan) bool {
	if ts.Length == 0 || other.Length == 0 {
		return false
	}
	if !ts.IsPairInOrder(other) {
		return other.OverlapsWith(ts)
	}
	return ts.Position <= other.Position && other.Position < ts.Position+ts.Length
}

func (ts *TextSpan) IsPairInOrder(other *TextSpan) bool {
	return ts.Position <= other.Position
}
