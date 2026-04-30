package helm

import (
	"fmt"
	"sync"
)

// MultiReleaseResult holds the result of fetching a single release.
type MultiReleaseResult struct {
	Release *Release
	Values  map[string]interface{}
	Err     error
	Name    string
}

// FetchMultipleReleases concurrently fetches release info and values for each
// of the provided release names. It returns one result per name, preserving
// the original order.
func FetchMultipleReleases(client Client, names []string) []MultiReleaseResult {
	results := make([]MultiReleaseResult, len(names))
	var wg sync.WaitGroup

	for i, name := range names {
		wg.Add(1)
		go func(idx int, relName string) {
			defer wg.Done()
			res := MultiReleaseResult{Name: relName}

			rel, err := client.GetRelease(relName)
			if err != nil {
				res.Err = fmt.Errorf("get release %q: %w", relName, err)
				results[idx] = res
				return
			}
			res.Release = rel

			vals, err := client.GetValues(relName)
			if err != nil {
				res.Err = fmt.Errorf("get values %q: %w", relName, err)
				results[idx] = res
				return
			}
			res.Values = vals
			results[idx] = res
		}(i, name)
	}

	wg.Wait()
	return results
}
