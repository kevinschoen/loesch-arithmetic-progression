package loesch

// IsLoesch reports whether n is a Loeschian number (OEIS A003136):
// n = x² + xy + y² for integers x, y. Equivalently, every prime p ≡ 2 (mod 3)
// appears with an even exponent in n's factorization. Zero is Loeschian.
func IsLoesch(n int64) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}

	for n%3 == 0 {
		n /= 3
	}

	for p := int64(2); p*p <= n; p++ {
		if n%p != 0 {
			continue
		}
		exp := 0
		for n%p == 0 {
			n /= p
			exp++
		}
		if p%3 == 2 && exp%2 != 0 {
			return false
		}
	}

	return n == 1 || n%3 != 2
}

// FirstN returns the first count Loeschian numbers in ascending order.
func FirstN(count int) []int64 {
	if count <= 0 {
		return nil
	}

	out := make([]int64, 0, count)
	for n := int64(0); len(out) < count; n++ {
		if IsLoesch(n) {
			out = append(out, n)
		}
	}
	return out
}
