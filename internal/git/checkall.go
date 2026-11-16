package git

import "sync"

// maxConcurrent caps how many git processes run at the same time.
const maxConcurrent = 8

// CheckAll checks every path concurrently and returns results in the same order as paths.
func CheckAll(paths []string) []RepoStatus {
	results := make([]RepoStatus, len(paths))
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			sem <- struct{}{}        // acquire a slot
			defer func() { <-sem }() // release it when done
			results[i] = Check(p)    // each goroutine writes only its own slot
		}(i, p)
	}

	wg.Wait()
	return results
}
