package token

import "bicep-go/util"

type TriviaType int

type Trivia struct {
	Type TriviaType
	Text string
	Span *util.TextSpan
}

const (
	WhitespaceTrivia TriviaType = iota
	SingleLineCommentTrivia
	MultiLineCommentTrivia
	NewLineTrivia
)

func NewTrivia(triviaType TriviaType, text string, span *util.TextSpan) *Trivia {
	return &Trivia{
		Type: triviaType,
		Text: text,
		Span: span,
	}
}
