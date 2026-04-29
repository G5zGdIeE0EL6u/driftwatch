package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
)

// Client wraps Helm action configuration for release operations.
type Client struct {
	cfg *action.Configuration
}

// NewClient creates a new Helm client for the given namespace.
func NewClient(namespace string, restConfig *rest.Config) (*Client, error) {
	cfg := new(action.Configuration)
	getter := newRESTClientGetter(namespace, restConfig)
	if err := cfg.Init(getter, namespace, "secret", func(format string, v ...interface{}) {
		// suppress helm debug logs
	}); err != nil {
		return nil, fmt.Errorf("initializing helm config: %w", err)
	}
	return &Client{cfg: cfg}, nil
}

// GetRelease retrieves a deployed Helm release by name.
func (c *Client) GetRelease(name string) (*release.Release, error) {
	get := action.NewGet(c.cfg)
	rel, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("getting release %q: %w", name, err)
	}
	return rel, nil
}

// GetReleaseValues returns the user-supplied values for a release.
func (c *Client) GetReleaseValues(name string) (map[string]interface{}, error) {
	rel, err := c.GetRelease(name)
	if err != nil {
		return nil, err
	}
	return rel.Config, nil
}

// GetRenderedManifests returns the rendered manifest string from a release.
func (c *Client) GetRenderedManifests(name string) (string, error) {
	rel, err := c.GetRelease(name)
	if err != nil {
		return "", err
	}
	return rel.Manifest, nil
}
