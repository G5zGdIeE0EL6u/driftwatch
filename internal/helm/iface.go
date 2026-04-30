package helm

import "helm.sh/helm/v3/pkg/release"

// ClientIface defines the operations used by higher-level packages
// so they can be easily mocked in tests.
type ClientIface interface {
	// GetRelease retrieves the named Helm release from the given namespace.
	GetRelease(name, namespace string) (*release.Release, error)

	// GetValues returns the user-supplied values for the named release.
	GetValues(name, namespace string) (map[string]interface{}, error)
}
