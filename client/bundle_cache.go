package keymaster

//BundleCache the type that represents the cache of bundles
type BundleCache struct {
	ApidURI string
	cache   []Bundle
	client  *ApidClient
}

//Bundle the type of the bundle to return
//TODO, not sure if this is even neccessary
type Bundle DeploymentBundle

//CreateBundleCache create the bundle cache
func CreateBundleCache(apidURL string, workingDirectory string) (*BundleCache, error) {

	client, err := CreateApidClient(apidURL, workingDirectory)

	if err != nil {
		return nil, err
	}

	bundleCache := &BundleCache{
		client: client,
	}

	return bundleCache, nil

}

//GetBundles returns all bundles via a channel.  Blocks until a change in the bundle set is detected
func (bundleCache *BundleCache) GetBundles() ([]Bundle, error) {

	return []Bundle{}, nil
}
