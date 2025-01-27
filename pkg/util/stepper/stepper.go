package stepper

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type Stepper[T any] struct {
	cur    T
	next   func(T) T
	locker sync.Mutex
}

func New[T any](init T, next func(T) T) *Stepper[T] {
	return &Stepper[T]{
		cur:  init,
		next: next,
	}
}

func (s *Stepper[T]) Next() T {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.cur = s.next(s.cur)

	return s.cur
}

func Int[T constraints.Integer]() *Stepper[T] {
	return New[T](0, func(t T) T {
		return t + 1
	})
}
