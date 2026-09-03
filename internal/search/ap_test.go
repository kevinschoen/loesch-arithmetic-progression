package search

import (
	"errors"
	"testing"
)

func TestExampleAP(t *testing.T) {
	if !IsAP(1, 6, 4) {
		t.Fatal("expected 1,7,13,19 to be a 4-term Loesch AP")
	}
}

func TestFindMinEndAPFourTerms(t *testing.T) {
	got, err := FindMinEndAP(4, 0, 100, nil)
	if err != nil {
		t.Fatalf("FindMinEndAP(4): %v", err)
	}
	if got.Start != 1 || got.Step != 6 || got.End != 19 {
		t.Fatalf("got start=%d step=%d end=%d, want start=1 step=6 end=19", got.Start, got.Step, got.End)
	}
}

func TestFindMinEndAPNotFound(t *testing.T) {
	_, err := FindMinEndAP(35, 0, 0, nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
