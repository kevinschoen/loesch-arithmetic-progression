package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kevinschoen/loesch-arithmetic-progression/internal/search"
)

func main() {
	terms := flag.Int("terms", 35, "number of terms in the arithmetic progression")
	minEnd := flag.Int64("min-end", 0, "start searching from this end value (resume a previous run)")
	maxEnd := flag.Int64("max-end", 50_000_000, "maximum end value to search")
	verbose := flag.Bool("verbose", false, "print progress while searching")
	flag.Parse()

	var progress func(int64, time.Duration)
	if *verbose {
		fmt.Fprintf(os.Stderr, "building sieve up to %d...\n", *maxEnd)
		progress = func(end int64, elapsed time.Duration) {
			fmt.Fprintf(os.Stderr, "searching end=%d elapsed=%s\n", end, elapsed.Round(time.Millisecond))
		}
	}

	result, err := search.FindMinEndAP(*terms, *minEnd, *maxEnd, progress)
	if err != nil {
		if errors.Is(err, search.ErrNotFound) {
			log.Fatalf("no %d-term progression found with end <= %d; try increasing -max-end", *terms, *maxEnd)
		}
		log.Fatal(err)
	}

	fmt.Printf("start=%d\n", result.Start)
	fmt.Printf("step=%d\n", result.Step)
	fmt.Printf("end=%d\n", result.End)
	fmt.Printf("terms=%d\n", result.Terms)
}
