// credit: https://github.com/ttdsuen/golang-stack/blob/main/stack.go

package util

type stackNode[T any] struct {
	data T
	next *stackNode[T]
}

type Stack[T any] struct {
	head *stackNode[T]
	Size int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		head: nil,
		Size: 0,
	}
}

func (s *Stack[T]) Push(data T) {
	newNode := &stackNode[T]{
		data: data,
		next: s.head,
	}
	s.head = newNode
	s.Size++
}

func (s *Stack[T]) Any() bool {
	return s.head != nil
}

func (s *Stack[T]) Peek() (T, bool) {
	var x T
	if !s.Any() {
		return x, false
	}
	return s.head.data, true
}

func (s *Stack[T]) Pop() (T, bool) {
	var x T
	if s.head == nil {
		return x, false
	}
	x = s.head.data
	s.head = s.head.next
	s.Size--
	return x, true
}
