package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("can push and pop 1 element", func(t *testing.T) {
		s := newStack[int](5)
		s.push(5)
		x, err := s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 5, x)
		}
	})
	t.Run("can push and pop 2 elements", func(t *testing.T) {
		s := newStack[int](5)
		s.push(5)
		s.push(3)
		x, err := s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 3, x)
		}
		x, err = s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 5, x)
		}
	})
	t.Run("should raise error when trying to pop from empty stack", func(t *testing.T) {
		s := newStack[int](5)
		_, err := s.pop()
		assert.ErrorIs(t, err, errStackEmpty)
	})
	t.Run("should only store the last x elements", func(t *testing.T) {
		s := newStack[int](3)
		s.push(1)
		s.push(2)
		s.push(3)
		s.push(4)
		assert.Len(t, s.l, 3)
		x, err := s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 4, x)
		}
		x, err = s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 3, x)
		}
		x, err = s.pop()
		if assert.NoError(t, err) {
			assert.Equal(t, 2, x)
		}
		_, err = s.pop()
		assert.ErrorIs(t, err, errStackEmpty)
	})
	t.Run("should panic if trying to init with invalid size", func(t *testing.T) {
		assert.Panics(t, func() {
			newStack[int](0)
		})
		assert.Panics(t, func() {
			newStack[int](-1)
		})
	})
}
