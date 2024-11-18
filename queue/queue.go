package queue

type Queue[T any] []T

func (q *Queue[T]) Push(x T) {
	*q = append(*q, x)
}

func (q *Queue[T]) Size() int {
	return len(*q)
}

func (q *Queue[T]) Pop() T {
	h := *q
	var el T
	l := len(h)
	el, *q = h[0], h[1:l]
	// Or use this instead for a Stack
	// el, *self = h[l-1], h[0:l-1]
	return el
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func NewQueueFromSlice[T any](slice []T) *Queue[T] {
	q := NewQueue[T]()

	// push all to the queue
	for _, item := range slice {
		q.Push(T(item))
	}

	return q
}
