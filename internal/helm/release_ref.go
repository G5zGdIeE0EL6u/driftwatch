package helm

import "fmt"

// ReleaseRef is a lightweight reference to a Helm release.
type ReleaseRef struct {
	Name      string
	Namespace string
}

// String returns a human-readable representation of the release reference.
func (r ReleaseRef) String() string {
	if r.Namespace == "" {
		return r.Name
	}
	return fmt.Sprintf("%s/%s", r.Namespace, r.Name)
}

// Key returns a unique string key for use in maps or caches.
func (r ReleaseRef) Key() string {
	return r.String()
}
