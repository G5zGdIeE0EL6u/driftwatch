package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
)

// GetValues returns the user-supplied values for the named Helm release in the
// given namespace. It merges chart defaults with user overrides, preferring
// user-supplied values so the result mirrors what was actually deployed.
func (c *Client) GetValues(namespace, name string) (map[string]interface{}, error) {
	cfg := new(action.Configuration)
	if err := cfg.Init(c.getter, namespace, "", func(_ string, _ ...interface{}) {}); err != nil {
		return nil, fmt.Errorf("helm config init: %w", err)
	}

	get := action.NewGetValues(cfg)
	get.AllValues = false // user-supplied values only

	vals, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("get values %s/%s: %w", namespace, name, err)
	}
	if vals == nil {
		return map[string]interface{}{}, nil
	}
	return vals, nil
}
