package queue

import (
	"errors"
	"shake/optional"
)

type Queue[T any] struct {
	items []T
	head  int
}

// Push adds an element to the queue.
func (q *Queue[T]) Push(x T) {
	if q.head > 0 {
		// Overwrite the position of the head if it's not at the top
		q.items[q.head-1] = x
		q.head-- // Move the head back one position
	} else {
		// Otherwise, append to the end
		q.items = append(q.items, x)
	}
}

// Size returns the number of elements currently in the queue.
func (q *Queue[T]) Size() int {
	return len(q.items) - q.head
}

// Pop removes and returns the element at the front of the queue.
func (q *Queue[T]) Pop() T {
	if q.Size() == 0 {
		var zero T
		return zero // Return the zero value if the queue is empty
	}
	el := q.items[q.head]
	q.head++
	return el
}

// GetSliceElement safely retrieves an element from a slice using negative or positive indexing.
func GetSliceElement[T any](slice []T, index int) (T, error) {
	if index < 0 {
		index = len(slice) + index
	}
	if index < 0 || index >= len(slice) {
		var zero T
		return zero, errors.New("index out of range")
	}
	return slice[index], nil
}

// Peek retrieves an element from the queue at a specific offset without removing it.
func (q *Queue[T]) Peek(offset int) optional.Optional[T] {
	if offset < -q.Size() || offset >= q.Size() {
		return optional.NewEmptyOptional[T]()
	}
	item, err := GetSliceElement(q.items[q.head:], offset)
	if err != nil {
		return optional.NewEmptyOptional[T]()
	}
	return optional.NewOptional(item)
}

// NewQueue creates a new empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
		head:  0,
	}
}

// NewQueueFromSlice creates a queue initialized with elements from a slice.
func NewQueueFromSlice[T any](slice []T) *Queue[T] {
	return &Queue[T]{
		items: slice,
		head:  0,
	}
}
