package queue

type Queue[T any] struct {
	ch chan T
}

func NewQueue[T any](size int) *Queue[T] {
	return &Queue[T]{ch: make(chan T, size)}
}

func (q *Queue[T]) Push(item T) {
	q.ch <- item
}

func (q *Queue[T]) Channel() chan T {
	return q.ch
}

func (q *Queue[T]) Close() {
	close(q.ch)
}
