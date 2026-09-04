package search

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
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

// isAP checks the AP membership, testing from both ends toward the middle so
// that failures in the interior (statistically harder) are caught sooner.
func isAP(start, step int64, terms int, set *loesch.Set) bool {
	lo, hi := 0, terms-1
	for lo <= hi {
		for _, k := range [2]int{lo, hi} {
			n := start + int64(k)*step
			if set != nil {
				if !set.Contains(n) {
					return false
				}
			} else {
				if !loesch.IsLoesch(n) {
					return false
				}
			}
		}
		lo++
		hi--
	}
	return true
}

// FindMinEndAP finds an arithmetic progression of terms Loesch numbers whose
// last term is as small as possible. Search scans end values in ascending order,
// starting from minEnd. Pass 0 to search from the beginning.
func FindMinEndAP(terms int, minEnd, maxEnd int64, progress func(end int64, elapsed time.Duration)) (Result, error) {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	return findMinEndAP(terms, minEnd, maxEnd, workers, progress)
}

func findMinEndAP(terms int, minEnd, maxEnd int64, workers int, progress func(end int64, elapsed time.Duration)) (Result, error) {
	if terms < 1 {
		return Result{}, fmt.Errorf("terms must be positive, got %d", terms)
	}
	if minEnd < 0 {
		return Result{}, fmt.Errorf("minEnd must be non-negative, got %d", minEnd)
	}
	if maxEnd < 0 {
		return Result{}, fmt.Errorf("maxEnd must be non-negative, got %d", maxEnd)
	}

	set := loesch.NewSet(maxEnd)

	startTime := time.Now()
	span := int64(terms - 1)

	// Distribute end values across workers. Each worker owns every
	// (workers)-th integer starting at (minEnd + workerID). Because we want
	// the globally minimal end, we interleave rather than block-partition:
	// worker 0 checks end=minEnd, minEnd+workers, …
	// worker 1 checks end=minEnd+1, minEnd+workers+1, …
	// …so all workers advance in lock-step and the first result found across
	// all workers is guaranteed to have the smallest possible end.

	type hit struct {
		result Result
		end    int64
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		best    hit
		hasBest bool
		// stopAt is the smallest end found so far; workers skip values >= stopAt.
		stopAt atomic.Int64
	)
	stopAt.Store(maxEnd + 1)

	// workerPos[w] holds the current end value for worker w.
	// The ticker reports the minimum across all workers — the true search frontier.
	workerPos := make([]atomic.Int64, workers)
	for w := 0; w < workers; w++ {
		workerPos[w].Store(minEnd + int64(w))
	}

	// Drive progress reports from a dedicated ticker goroutine so they are
	// never blocked by tight worker loops.
	done := make(chan struct{})
	if progress != nil {
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					min := workerPos[0].Load()
					for w := 1; w < workers; w++ {
						if v := workerPos[w].Load(); v < min {
							min = v
						}
					}
					progress(min, time.Since(startTime))
				}
			}
		}()
	}

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w
		go func() {
			defer wg.Done()
			for end := minEnd + int64(w); end <= maxEnd; end += int64(workers) {
				if end >= stopAt.Load() {
					return
				}

				workerPos[w].Store(end)

				if !set.Contains(end) {
					continue
				}
				if end < span {
					continue
				}

				// Try all steps d such that start = end - span*d is Loeschian.
				// Walk start downward in steps of d=1,2,… but only visit values
				// where start itself is Loeschian — O(end/span) per end.
				for start := end - span; start >= 0; start -= span {
					if end >= stopAt.Load() {
						break
					}
					if !set.Contains(start) {
						continue
					}
					step := (end - start) / span
					if isAP(start, step, terms, set) {
						mu.Lock()
						if !hasBest || end < best.end {
							best = hit{
								result: Result{Start: start, Step: step, End: end, Terms: terms},
								end:    end,
							}
							hasBest = true
							stopAt.Store(end) // no need to check any end >= this
						}
						mu.Unlock()
						break
					}
				}
			}
		}()
	}
	wg.Wait()
	close(done)

	if !hasBest {
		return Result{}, fmt.Errorf("%w (max end %d)", ErrNotFound, maxEnd)
	}
	return best.result, nil
}
