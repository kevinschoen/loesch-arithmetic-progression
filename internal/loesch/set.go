package loesch

// Set is a fast membership set for Loeschian numbers up to max.
type Set struct {
	max   int64
	loesch []bool
}

// NewSet precomputes all Loeschian numbers from 0 through max inclusive.
func NewSet(max int64) *Set {
	s := &Set{
		max:    max,
		loesch: make([]bool, max+1),
	}
	for n := int64(0); n <= max; n++ {
		s.loesch[n] = IsLoesch(n)
	}
	return s
}

// Contains reports whether n is in the precomputed Loeschian set.
func (s *Set) Contains(n int64) bool {
	return n >= 0 && n <= s.max && s.loesch[n]
}
