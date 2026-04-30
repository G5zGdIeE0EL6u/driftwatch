package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
)

// ValuesSource indicates which set of values to retrieve from a release.
type ValuesSource int

const (
	// UserSupplied returns only the values explicitly provided by the user.
	UserSupplied ValuesSource = iota
	// ComputedAll returns the merged set of default + user-supplied values.
	ComputedAll
)

// GetValues retrieves the Helm values for the named release in the given
// namespace. Use src to choose between user-supplied or fully-computed values.
func (c *Client) GetValues(release, namespace string, src ValuesSource) (map[string]interface{}, error) {
	cfg := new(action.Configuration)
	if err := cfg.Init(c.getter, namespace, "", func(_ string, _ ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("helm config init: %w", err)
	}

	get := action.NewGetValues(cfg)

	switch src {
	case ComputedAll:
		get.AllValues = true
	default:
		get.AllValues = false
	}

	vals, err := get.Run(release)
	if err != nil {
		return nil, fmt.Errorf("get values %q: %w", release, err)
	}

	if vals == nil {
		return map[string]interface{}{}, nil
	}

	return vals, nil
}
