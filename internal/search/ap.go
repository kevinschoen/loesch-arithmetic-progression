package search

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/kevinschoen/loesch-arithmetic-progression/internal/loesch"
)

var ErrNotFound = errors.New("no arithmetic progression found within max end")

type Result struct {
	Start int64
	Step  int64
	End   int64
	Terms int
}

// IsAP reports whether start + k*step for k = 0..terms-1 are all Loeschian.
func IsAP(start, step int64, terms int) bool {
	return isAP(start, step, terms, nil)
}

func isAP(start, step int64, terms int, set *loesch.Set) bool {
	for k := range terms {
		n := start + int64(k)*step
		if set != nil {
			if !set.Contains(n) {
				return false
			}
			continue
		}
		if !loesch.IsLoesch(n) {
			return false
		}
	}
	return true
}

// FindMinEndAP finds an arithmetic progression of terms Loesch numbers whose
// last term is as small as possible. Search scans end values in ascending order.
func FindMinEndAP(terms int, maxEnd int64, progress func(end int64, elapsed time.Duration)) (Result, error) {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	return findMinEndAP(terms, maxEnd, workers, progress)
}

func findMinEndAP(terms int, maxEnd int64, workers int, progress func(end int64, elapsed time.Duration)) (Result, error) {
	if terms < 1 {
		return Result{}, fmt.Errorf("terms must be positive, got %d", terms)
	}
	if maxEnd < 0 {
		return Result{}, fmt.Errorf("maxEnd must be non-negative, got %d", maxEnd)
	}

	set := loesch.NewSet(maxEnd)

	startTime := time.Now()
	lastReport := startTime
	span := int64(terms - 1)

	for end := int64(0); end <= maxEnd; end++ {
		if progress != nil && time.Since(lastReport) >= time.Second {
			progress(end, time.Since(startTime))
			lastReport = time.Now()
		}

		if !set.Contains(end) {
			continue
		}

		maxStep := end / span
		if maxStep == 0 {
			continue
		}

		if result, ok := searchSteps(end, span, terms, maxStep, set, workers); ok {
			return result, nil
		}
	}

	return Result{}, fmt.Errorf("%w (max end %d)", ErrNotFound, maxEnd)
}

func searchSteps(end, span int64, terms int, maxStep int64, set *loesch.Set, workers int) (Result, bool) {
	if maxStep <= 0 {
		return Result{}, false
	}

	type hit struct {
		start int64
		step  int64
	}

	jobs := make(chan int64, workers*2)
	var wg sync.WaitGroup
	var once sync.Once
	var found hit
	done := make(chan struct{})

	worker := func() {
		defer wg.Done()
		for step := range jobs {
			select {
			case <-done:
				return
			default:
			}

			start := end - span*step
			if isAP(start, step, terms, set) {
				once.Do(func() {
					found = hit{start: start, step: step}
					close(done)
				})
				return
			}
		}
	}

	wg.Add(workers)
	for range workers {
		go worker()
	}

loop:
	for step := int64(1); step <= maxStep; step++ {
		select {
		case <-done:
			break loop
		default:
			start := end - span*step
			if set.Contains(start) {
				jobs <- step
			}
		}
	}
	close(jobs)
	wg.Wait()

	if found.step == 0 {
		return Result{}, false
	}

	return Result{
		Start: found.start,
		Step:  found.step,
		End:   end,
		Terms: terms,
	}, true
}
