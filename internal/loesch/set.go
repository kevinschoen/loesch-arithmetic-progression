package loesch

import "sort"

// Set is a fast membership set for Loeschian numbers up to max.
type Set struct {
	max    int64
	loesch []bool
	sorted []int64 // all Loeschian numbers in ascending order
}

// NewSet precomputes all Loeschian numbers from 0 through max inclusive.
// It uses a generative sieve: enumerate all x²+xy+y² pairs instead of
// testing each integer individually, reducing build cost from O(n√n) to O(n).
func NewSet(max int64) *Set {
	s := &Set{
		max:    max,
		loesch: make([]bool, max+1),
	}
	for x := int64(0); ; x++ {
		v0 := x * x // x²+x*0+0² = x²
		if v0 > max {
			break
		}
		for y := int64(0); y <= x; y++ {
			v := x*x + x*y + y*y
			if v > max {
				break
			}
			s.loesch[v] = true
		}
	}
	// Build sorted slice from the bitmap.
	for n := int64(0); n <= max; n++ {
		if s.loesch[n] {
			s.sorted = append(s.sorted, n)
		}
	}
	return s
}

// Contains reports whether n is in the precomputed Loeschian set.
func (s *Set) Contains(n int64) bool {
	return n >= 0 && n <= s.max && s.loesch[n]
}

// Below returns a slice of all Loeschian numbers strictly less than n,
// in ascending order.
func (s *Set) Below(n int64) []int64 {
	i := sort.Search(len(s.sorted), func(i int) bool { return s.sorted[i] >= n })
	return s.sorted[:i]
}
