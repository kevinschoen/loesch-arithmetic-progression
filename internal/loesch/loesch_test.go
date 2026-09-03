package loesch

import (
	"slices"
	"testing"
)

func TestFirstNLoeschNumbers(t *testing.T) {
	want := []int64{0, 1, 3, 4, 7, 9, 12, 13, 16, 19, 21, 25}
	got := FirstN(len(want))
	if !slices.Equal(got, want) {
		t.Fatalf("FirstN(%d) = %v, want %v", len(want), got, want)
	}
}

func TestNonLoeschNumbers(t *testing.T) {
	nonLoesch := []int64{2, 5, 6, 8, 10, 11, 14, 15, 17, 18, 20, 22, 23, 26}
	for _, n := range nonLoesch {
		if IsLoesch(n) {
			t.Fatalf("%d should not be Loeschian", n)
		}
	}
}

func TestNegativeNotLoesch(t *testing.T) {
	if IsLoesch(-1) {
		t.Fatal("-1 should not be Loeschian")
	}
}

func TestPrimeTwoModThreeEvenPower(t *testing.T) {
	// 7 is Loeschian; 7*7 = 49 is Loeschian; 7 alone in factorization with odd power is not.
	if !IsLoesch(49) {
		t.Fatal("49 = 7² should be Loeschian")
	}
	if IsLoesch(14) {
		t.Fatal("14 = 2*7 should not be Loeschian")
	}
}
