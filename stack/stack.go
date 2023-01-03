package stack

import "errors"

type Stack[A any] []A

func New[A any]() Stack[A] {
	return make([]A, 0)
}

var ErrEmpty = errors.New("Stack Underflow")

func (s *Stack[A]) Pop() (A, error) {
	l := len(*s)
	if l == 0 {
		var zero A
		return zero, ErrEmpty
	}
	v := (*s)[l-1]
	*s = (*s)[:l-1]
	return v, nil
}

func (s *Stack[A]) MustPop() A {
	l := len(*s)
	v := (*s)[l-1]
	*s = (*s)[:l]
	return v
}

func (s *Stack[A]) Push(v A) {
	*s = append(*s, v)
}

// ForEach in order, top first
func (s Stack[A]) ForEach(f func(A)) {
	l := len(s)
	for i := range s {
		f(s[l-i-1])
	}
}
