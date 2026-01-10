package functools

import (
	"errors"
	"strings"
)

func Must[T any](a T, e error) T {
	if e != nil {
		panic(e)
	}
	return a
}

func MapWithError[T, V any](f func(T) (V, error), xs []T) ([]V, error) {
	result := make([]V, len(xs))
	var encounteredErrors []string
	for i, t := range xs {
		v, err := f(t)
		if err != nil {
			encounteredErrors = append(encounteredErrors, err.Error())
			continue
		}
		result[i] = v
	}

	if len(encounteredErrors) > 0 {
		return nil, errors.New(strings.Join(encounteredErrors, "; "))
	}

	return result, nil
}

func MapShortCircuit[T, V any](f func(T) (V, error), xs []T) ([]V, error) {
	result := make([]V, len(xs))

	for i, t := range xs {
		v, err := f(t)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}

	return result, nil
}

type Pair[T, V any] struct {
	First  T
	Second V
}

func Zip[T, V any](xs []T, ys []V) []Pair[T, V] {
	var pairs []Pair[T, V]
	min_length := len(xs)
	if len(ys) < len(xs) {
		min_length = len(ys)
	}
	for i := 0; i < min_length; i += 1 {
		pairs = append(pairs, Pair[T, V]{xs[i], ys[i]})
	}
	return pairs
}

func Reduce[T, V any](f func(V, T) V, initial V, xs []T) V {
	accumulator := initial
	for _, x := range xs {
		accumulator = f(accumulator, x)
	}
	return accumulator
}

func All[T any](f func(T) bool, xs []T) bool {
	for _, x := range xs {
		if f(x) == false {
			return false
		}
	}
	return true
}

func Any[T any](f func(T) bool, xs []T) bool {
	for _, x := range xs {
		if f(x) == true {
			return true
		}
	}
	return false
}
