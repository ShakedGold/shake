package queue

import (
	"encoding/json"
	"fmt"
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
func (q *Queue[T]) Pop() *T {
	if q.Size() == 0 {
		return nil // Return the zero value if the queue is empty
	}
	el := &q.items[q.head]
	q.head++
	return el
}

func (q *Queue[T]) TryPop() (*T, error) {
	_, err := q.Peek(0)
	if err != nil {
		return nil, err
	}
	item := q.Pop()
	return item, nil
}

// Peek retrieves an element from the queue at a specific offset without removing it.
func (q *Queue[T]) Peek(offset int) (*T, error) {
	if offset < -q.Size() || offset >= q.Size() {
		return nil, fmt.Errorf("offset: %d not in range: %d", offset, q.Size())
	}
	return &q.items[offset+q.head], nil
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

// Custom MarshalJSON method for TokenType
func (q *Queue[T]) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(q.items, "", "\t")
}
