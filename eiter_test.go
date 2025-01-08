// Copyright 2025 Bob Vawter (bob@vawter.org)
// SPDX-License-Identifier: MIT

package eiter_test

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"vawter.tech/eiter"
)

// Counter emits the requested number of values and then an error.
func Counter(ct int) eiter.Seq[int] {
	return eiter.Of(func(yield func(int) bool) error {
		for i := range ct {
			if !yield(i) {
				return nil
			}
		}
		return errors.New("Error World!")
	})
}

// Counter2 emits the requested number of values and then an error.
func Counter2(ct int) eiter.Seq2[int, string] {
	return eiter.Of2(func(yield func(int, string) bool) error {
		for i := range ct {
			if !yield(i, strconv.Itoa(i)) {
				return nil
			}
		}
		return errors.New("Error World!")
	})
}

func ExampleJust() {
	it := eiter.Just(slices.Values([]int{1, 2, 3}))
	for entry := range it {
		if err := entry.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(entry.Value())
		}
	}
	// Output:
	// 1
	// 2
	// 3
}

func ExampleJust2() {
	it := eiter.Just2(slices.All([]int{1, 2, 3}))
	for entry := range it {
		if err := entry.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(entry.Key(), entry.Value())
		}
	}
	// Output:
	// 0 1
	// 1 2
	// 2 3
}

func ExampleOf() {
	for entry := range Counter(3) {
		if err := entry.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(entry.Value())
		}
	}
	// Output:
	// 0
	// 1
	// 2
	// Error World!
}

func ExampleOf2() {
	for entry := range Counter2(3) {
		if err := entry.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(entry.Key(), entry.Value())
		}
	}
	// Output:
	// 0 0
	// 1 1
	// 2 2
	// Error World!
}

func SomeOtherCallbackAPI(fn func() error) error {
	return fn()
}

// The [eiter.ErrStop] sentinel value can be used if the generator is
// constructed over some other callback-based API.
func ExampleOf_stop() {
	it := eiter.Of(func(yield func(int) bool) error {
		count := 0
		return SomeOtherCallbackAPI(func() error {
			yield(count)
			count++
			return eiter.ErrStop
		})
	})

	for r := range it {
		if err := r.Err(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r.Value())
		}
	}
	// Output:
	// 0
}

func ExampleSeq_Unwrap() {
	var err error
	for i := range Counter(3).Unwrap(&err) {
		fmt.Println(i)
	}
	// This is similar to checking sql.Rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 0
	// 1
	// 2
	// Error World!
}

func ExampleSeq_Unwrap_collect() {
	var err error
	values := slices.Collect(Counter(3).Unwrap(&err))
	fmt.Println(values)
	// This is similar to checking sql.Rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// [0 1 2]
	// Error World!
}

func ExampleSeq2_Unwrap() {
	var err error
	for k, v := range Counter2(3).Unwrap(&err) {
		fmt.Println(k, v)
	}
	// This is similar to checking sql.Rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 0 0
	// 1 1
	// 2 2
	// Error World!
}

func ExampleSeq2_Unwrap_collect() {
	var err error
	values := maps.Collect(Counter2(3).Unwrap(&err))
	fmt.Println(values)
	// This is similar to checking sql.Rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// map[0:0 1:1 2:2]
	// Error World!
}
