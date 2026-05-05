package helm

import (
	"context"
	"fmt"

	helm "helm.sh/helm/v3/pkg/action"
)

// NamespaceLister lists Helm releases across one or more namespaces.
type NamespaceLister struct {
	cfgFn func(namespace string) (*helm.Configuration, error)
}

// NewNamespaceLister creates a NamespaceLister using the provided kubeconfig path.
func NewNamespaceLister(kubeconfig, context string) (*NamespaceLister, error) {
	return &NamespaceLister{
		cfgFn: func(namespace string) (*helm.Configuration, error) {
			cfg := new(helm.Configuration)
			getter := newRESTClientGetter(namespace, kubeconfig, context)
			if err := cfg.Init(getter, namespace, "secret", func(format string, v ...interface{}) {}); err != nil {
				return nil, fmt.Errorf("init helm config for namespace %q: %w", namespace, err)
			}
			return cfg, nil
		},
	}, nil
}

// ListReleaseNames returns all release names found in the given namespaces.
// If namespaces is empty, it queries the "" (all) namespace.
func (l *NamespaceLister) ListReleaseNames(ctx context.Context, namespaces []string) ([]ReleaseRef, error) {
	if len(namespaces) == 0 {
		namespaces = []string{""}
	}

	var results []ReleaseRef
	for _, ns := range namespaces {
		refs, err := l.listInNamespace(ctx, ns)
		if err != nil {
			return nil, err
		}
		results = append(results, refs...)
	}
	return results, nil
}

func (l *NamespaceLister) listInNamespace(_ context.Context, namespace string) ([]ReleaseRef, error) {
	cfg, err := l.cfgFn(namespace)
	if err != nil {
		return nil, err
	}

	list := helm.NewList(cfg)
	list.All = true
	list.AllNamespaces = namespace == ""

	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("list releases in namespace %q: %w", namespace, err)
	}

	refs := make([]ReleaseRef, 0, len(releases))
	for _, r := range releases {
		refs = append(refs, ReleaseRef{Name: r.Name, Namespace: r.Namespace})
	}
	return refs, nil
}
