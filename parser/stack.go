package parser

import (
	"errors"
	"slices"
)

var errStackEmpty = errors.New("stack is empty")

// stack is a generic stack with limited size
type stack[T any] struct {
	l    []T
	size int
}

// newStack returns a new stack
func newStack[T any](size int) stack[T] {
	if size < 1 {
		panic("size must be 1 or larger")
	}
	s := stack[T]{
		l:    make([]T, 0, size),
		size: size,
	}
	return s
}

// push pushed a new element on the stack.
// If the stack reached it's limit it will discard the oldest element.
func (s *stack[T]) push(v T) {
	if len(s.l) == s.size {
		s.l = slices.Delete(s.l, 0, 1)
	}
	s.l = append(s.l, v)
}

// pop removes and returns the newest element from the stack.
func (s *stack[T]) pop() (T, error) {
	if len(s.l) == 0 {
		var x T
		return x, errStackEmpty
	}
	last := len(s.l) - 1
	v := s.l[last]
	s.l = slices.Delete(s.l, last, last+1)
	return v, nil
}

// isEmpty reports wether the stack is empty
func (s *stack[T]) isEmpty() bool {
	return len(s.l) == 0
}
