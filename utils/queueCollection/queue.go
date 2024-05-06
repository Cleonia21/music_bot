package queueCollection

import "errors"

const queueMaxLen = 20

type queue[T any] struct {
	len   int
	queue chan T
}

func newQueue[T any]() *queue[T] {
	q := queue[T]{
		len:   0,
		queue: make(chan T, 20),
	}
	return &q
}

func (q *queue[T]) set(elem T) error {
	if q.len == queueMaxLen {
		return errors.New("exceeding the maximum queue length")
	}
	q.len++
	q.queue <- elem
	return nil
}

func (q *queue[T]) get() (elem T, err error) {
	if q.len == 0 {
		return elem, errors.New("no items have been added to the queue")
	}
	elem = <-q.queue
	q.len--
	return elem, nil
}
