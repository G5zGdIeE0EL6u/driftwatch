package helm

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"helm.sh/helm/v3/pkg/storage"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/restmapper"
)

// fakeRESTClientGetter satisfies helm's genericclioptions.RESTClientGetter
// interface using an in-memory Helm storage backend, enabling unit tests
// without a real cluster.
type fakeRESTClientGetter struct {
	store *storage.Storage
}

func (f *fakeRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return &rest.Config{Host: "http://localhost"}, nil
}

func (f *fakeRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return nil, nil
}

func (f *fakeRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	groupResources := []*restmapper.APIGroupResources{}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	return mapper, nil
}

func (f *fakeRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return clientcmd.NewDefaultClientConfig(
		clientcmd.NewDefaultPathOptions().GetDefaultFilename,
		&clientcmd.ConfigOverrides{},
	)
}
