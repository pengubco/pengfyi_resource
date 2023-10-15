package runners

import "errors"

// RoundRobinContainer returns elements one by one, in a round robin manner.
type RoundRobinContainer[T any] struct {
	elements  []T
	size      int
	nextIndex int
}

// NewRoundRobinContainer creates a new container that round robins over the given elements.
func NewRoundRobinContainer[T any](elements ...T) (*RoundRobinContainer[T], error) {
	if len(elements) == 0 {
		return nil, errors.New("elements cannot be empty")
	}
	n := len(elements)
	newElements := make([]T, n)
	copy(newElements, elements)
	return &RoundRobinContainer[T]{
		elements: newElements,
		size:     n,
	}, nil
}

// Next returns the next element in a round robin manner.
func (c *RoundRobinContainer[T]) Next() T {
	ret := c.elements[c.nextIndex]
	c.nextIndex++
	if c.nextIndex == c.size {
		c.nextIndex = 0
	}
	return ret
}
