package helm

import (
	"context"
	"fmt"
	"sync"
)

// NamespaceDriftTarget groups a release reference with its fetched release data.
type NamespaceDriftTarget struct {
	Ref     ReleaseRef
	Release *Release
	Err     error
}

// FetchAllNamespaceReleases concurrently fetches all releases across the
// provided namespaces using the given HelmClient. Results preserve the order
// of the input refs slice.
func FetchAllNamespaceReleases(ctx context.Context, client HelmClient, lister *NamespaceLister, namespaces []string) ([]NamespaceDriftTarget, error) {
	if len(namespaces) == 0 {
		namespaces = []string{"default"}
	}

	refs, err := lister.ListReleaseNames(ctx, namespaces)
	if err != nil {
		return nil, fmt.Errorf("listing release names: %w", err)
	}

	results := make([]NamespaceDriftTarget, len(refs))
	var wg sync.WaitGroup

	for i, ref := range refs {
		wg.Add(1)
		go func(idx int, r ReleaseRef) {
			defer wg.Done()
			rel, fetchErr := client.GetRelease(ctx, r.Namespace, r.Name)
			results[idx] = NamespaceDriftTarget{
				Ref:     r,
				Release: rel,
				Err:     fetchErr,
			}
		}(i, ref)
	}

	wg.Wait()
	return results, nil
}
