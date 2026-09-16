package loesch

// Set is a fast membership set for Loeschian numbers up to max.
type Set struct {
	max     int64
	loesch  []bool
	Buckets [][]int64 // Buckets[residue] = sorted Loeschian numbers ≡ residue (mod stride)
	Stride  int64
}

// NewSet precomputes all Loeschian numbers from 0 through max inclusive,
// optionally organizing them into buckets by residue modulo stride.
// It uses a generative sieve: enumerate all x²+xy+y² pairs instead of
// testing each integer individually, reducing build cost from O(n√n) to O(n).
// If stride > 0, buckets are populated for efficient residue-class iteration.
func NewSet(max, stride int64) *Set {
	s := &Set{
		max:    max,
		loesch: make([]bool, max+1),
		Stride: stride,
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

	// Populate buckets if stride is provided
	if stride > 0 {
		s.Buckets = make([][]int64, stride)
		for n := int64(0); n <= max; n++ {
			if s.loesch[n] {
				residue := n % stride
				s.Buckets[residue] = append(s.Buckets[residue], n)
			}
		}
	}

	return s
}

// Contains reports whether n is in the precomputed Loeschian set.
func (s *Set) Contains(n int64) bool {
	return n >= 0 && n <= s.max && s.loesch[n]
}
