// Copyright 2025 Bob Vawter (bob@vawter.org)
// SPDX-License-Identifier: MIT

// Package eiter supports iterators that may encounter errors.
package eiter

import (
	"errors"
	"iter"
)

// ErrStop is a sentinel value recognized by [Of] and [Of2].
var ErrStop = errors.New("stop")

// An Entry contains either a value or an error.
type Entry[T any] struct {
	err   error
	value T
}

func (r *Entry[T]) Err() error { return r.err }
func (r *Entry[T]) Value() T   { return r.value }

// An Entry2 contains either a key/value pair or an error.
type Entry2[K, V any] struct {
	err   error
	key   K
	value V
}

func (r *Entry2[K, V]) Err() error { return r.err }
func (r *Entry2[K, V]) Key() K     { return r.key }
func (r *Entry2[K, V]) Value() V   { return r.value }

// A Seq is a sequence that may report an error.
type Seq[T any] iter.Seq[*Entry[T]]

// Unwrap converts a fallible sequence to a stdlib sequence. Iteration
// will stop on the first error, which will be stored into the given
// pointer.
func (s Seq[T]) Unwrap(err *error) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := range s {
			if e.Err() != nil {
				*err = e.Err()
				return
			}
			if !yield(e.Value()) {
				return
			}
		}
	}
}

// A Seq2 is a sequence of key/value pairs that may report an error.
type Seq2[K, V any] iter.Seq[*Entry2[K, V]]

// Unwrap converts a fallible sequence to a stdlib sequence. Iteration
// will stop on the first error, which will be stored into the given
// pointer.
func (s Seq2[K, V]) Unwrap(err *error) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range s {
			if e.Err() != nil {
				*err = e.Err()
				return
			}
			if !yield(e.Key(), e.Value()) {
				return
			}
		}
	}
}

// Just wraps an existing sequence.
func Just[T any](it iter.Seq[T]) Seq[T] {
	return func(yield func(*Entry[T]) bool) {
		for e := range it {
			if !yield(&Entry[T]{value: e}) {
				return
			}
		}
	}
}

// Just2 wraps an existing sequence.
func Just2[K, V any](it iter.Seq2[K, V]) Seq2[K, V] {
	return func(yield func(*Entry2[K, V]) bool) {
		for k, v := range it {
			if !yield(&Entry2[K, V]{key: k, value: v}) {
				return
			}
		}
	}
}

// Of constructs a fallible iterator from a generator function that may
// return an error. The generator function may either return nil or
// [ErrStop] to indicate successful termination.
func Of[T any](generator func(yield func(T) bool) error) Seq[T] {
	return func(yield func(*Entry[T]) bool) {
		if err := generator(func(t T) bool {
			return yield(&Entry[T]{value: t})
		}); err != nil {
			if !errors.Is(err, ErrStop) {
				yield(&Entry[T]{err: err})
			}
			return
		}
	}
}

// Of2 constructs a fallible iterator from a generator function that
// may return an error.
func Of2[K, V any](generator func(yield func(K, V) bool) error) Seq2[K, V] {
	return func(yield func(*Entry2[K, V]) bool) {
		if err := generator(func(k K, v V) bool {
			return yield(&Entry2[K, V]{key: k, value: v})
		}); err != nil {
			if !errors.Is(err, ErrStop) {
				yield(&Entry2[K, V]{err: err})
			}
			return
		}
	}
}
