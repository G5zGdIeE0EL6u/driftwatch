package helm

import (
	"fmt"
	"time"

	"github.com/yourorg/driftwatch/internal/cache"
	helm "helm.sh/helm/v3/pkg/release"
)

// CachedClient wraps a ReleaseClient and caches results to reduce
// repeated calls to the Kubernetes API during a single drift-check run.
type CachedClient struct {
	inner  *Client
	cache  *cache.Cache
}

// NewCachedClient creates a CachedClient with the supplied TTL.
func NewCachedClient(inner *Client, ttl time.Duration) *CachedClient {
	return &CachedClient{
		inner: inner,
		cache: cache.New(ttl),
	}
}

// GetRelease returns a Helm release, using the cache when possible.
func (c *CachedClient) GetRelease(namespace, name string) (*helm.Release, error) {
	key := fmt.Sprintf("release::%s::%s", namespace, name)
	if v, ok := c.cache.Get(key); ok {
		return v.(*helm.Release), nil
	}
	rel, err := c.inner.GetRelease(namespace, name)
	if err != nil {
		return nil, err
	}
	c.cache.Set(key, rel)
	return rel, nil
}

// GetValues returns chart values for a release, using the cache when possible.
func (c *CachedClient) GetValues(namespace, name string) (map[string]interface{}, error) {
	key := fmt.Sprintf("values::%s::%s", namespace, name)
	if v, ok := c.cache.Get(key); ok {
		return v.(map[string]interface{}), nil
	}
	vals, err := c.inner.GetValues(namespace, name)
	if err != nil {
		return nil, err
	}
	c.cache.Set(key, vals)
	return vals, nil
}

// Invalidate removes cached entries for the given release.
func (c *CachedClient) Invalidate(namespace, name string) {
	c.cache.Delete(fmt.Sprintf("release::%s::%s", namespace, name))
	c.cache.Delete(fmt.Sprintf("values::%s::%s", namespace, name))
}
