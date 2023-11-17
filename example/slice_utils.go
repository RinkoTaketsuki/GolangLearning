package example

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
)

func SprintSlice[S ~[]E, E any](s S) string {
	return fmt.Sprintf("len: %d\tcap: %d\t0addr:%p\t%v\n", len(s), cap(s), &s[0], s)
}

func FprintSlice[S ~[]E, E any](w io.Writer, s S) (int, error) {
	return fmt.Fprint(w, SprintSlice(s))
}

func PrintSlice[S ~[]E, E any](s S) {
	FprintSlice(os.Stdout, s)
}

func Clone[S ~[]E, E any](s S) S {
	ret := make(S, len(s))
	copy(ret, s)
	return ret
}

func PopBack[S ~[]E, E any](s S) (S, E) {
	return s[:len(s)-1], s[len(s)-1]
}

func PushFront[S ~[]E, E any](s S, e ...E) S {
	return append(e, s...)
}

func PopFront[S ~[]E, E any](s S) (S, E) {
	return s[1:], s[0]
}

func Insert[S ~[]E, E any](s S, i int, e E) S {
	if i < 0 || i > len(s) {
		panic("invalid args of InsertMany")
	}
	var zeroValue E
	s = append(s, zeroValue) // avoid memory leak
	copy(s[i+1:], s[i:])
	s[i] = e
	return s
}

func InsertMany[S ~[]E, E any](s S, i int, e ...E) S {
	if i < 0 || i > len(s) {
		panic("invalid args of InsertMany")
	}
	if n := len(s) + len(e); n <= cap(s) {
		s2 := s[:n]
		copy(s2[i+len(e):], s[i:])
		copy(s2[i:], e)
		return s2
	}
	s2 := make(S, len(s)+len(e))
	copy(s2, s[:i])
	copy(s2[i:], e)
	copy(s2[i+len(e):], s[i:])
	return s2
}

func Delete[S ~[]E, E any](s S, i int) S {
	if i < 0 || i >= len(s) {
		panic("invalid args of Delete")
	}
	copy(s[i:], s[i+1:])
	var zeroValue E
	s[len(s)-1] = zeroValue // avoid memory leak
	return s[:len(s)-1]
}

func Filter[S ~[]E, E any](s S, f func(E) bool) S {
	n := 0
	for i := range s {
		if f(s[i]) {
			s[n] = s[i]
			n++
		}
	}
	return s[:n]
}

// Delete elements in s[i:j]
func Cut[S ~[]E, E any](s S, i, j int) S {
	if i < 0 || i > j || j > len(s) {
		panic("invalid args of Cut")
	}
	copy(s[i:], s[j:])
	var zeroValue E
	for k := len(s) - j + i; k < len(s); k++ {
		s[k] = zeroValue // avoid memory leak
	}
	return s[:len(s)-j+i]
}

func Reverse[S ~[]E, E any](s S) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}
}

func Batch[S ~[]E, E any](s S, sz int) []S {
	if sz <= 0 {
		panic("invalid args of Expand")
	}
	if sz >= len(s) {
		return []S{s}
	}
	batches := make([]S, 0, (len(s)+sz-1)/sz)
	for sz < len(s) {
		s, batches = s[sz:], append(batches, s[0:sz:sz])
	}
	return append(batches, s)
}

// Insert n zero-value elements at position i:
func Expand[S ~[]E, E any](s S, i, n int) S {
	if n < 0 || i < 0 || i > len(s) {
		panic("invalid args of Expand")
	}
	return append(s[:i], append(make(S, n), s[i:]...)...)
}

func Shuffle[S ~[]E, E any](s S) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

func SortAndDeduplicate[S ~[]E, Lt func(E, E) bool, E any](s S, lt Lt) S {
	sort.Slice(s, func(i, j int) bool {
		return lt(s[i], s[j])
	})
	eq := func(e1, e2 E) bool {
		return !(lt(e1, e2) || lt(e2, e1))
	}
	j := 0
	for i := 1; i < len(s); i++ {
		if eq(s[i], s[j]) {
			continue
		}
		j++
		s[j] = s[i]
	}
	return s[:j+1]
}
