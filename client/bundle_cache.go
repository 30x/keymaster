package client

import "fmt"

//BundleCache the type that represents the cache of bundles
type BundleCache struct {
	lastEtag string
	timeout  int
	cache    map[string]*DeploymentBundle
	client   *ApidClient
}

// //Bundle the type of the bundle to return
// //TODO, not sure if this is even neccessary
// type Bundle

//CreateBundleCache create the bundle cache
func CreateBundleCache(apidURL string, workingDirectory string, timeout int) (*BundleCache, error) {

	client, err := CreateApidClient(apidURL, workingDirectory)

	if err != nil {
		return nil, err
	}

	bundleCache := &BundleCache{
		client:  client,
		timeout: timeout,
		cache:   make(map[string]*DeploymentBundle),
	}

	return bundleCache, nil

}

//GetBundles returns all bundles via a channel.  Blocks until a change in the bundle set is detected
func (bundleCache *BundleCache) GetBundles() ([]*DeploymentBundle, error) {

	response, err := bundleCache.client.PollDeployments(bundleCache.lastEtag, bundleCache.timeout)

	if err != nil {
		return nil, err
	}

	//iterate through the response, and check the cache for existing ids

	bundles := []*DeploymentBundle{}

	for _, bundleEntry := range response.Bundles {

		bundle, ok := bundleCache.cache[bundleEntry.BundleID]

		//it's not local get it from the server
		if !ok {

			bundle, err = bundleCache.client.GetBundle(bundleEntry)

			if bundle == nil {
				return nil, fmt.Errorf("Bundle %+v could not be found, and was returned in the deployment", bundleEntry)
			}
			if err != nil {
				return nil, err
			}

			bundleCache.cache[bundle.BundleID] = bundle

		}

		bundles = append(bundles, bundle)

	}

	//TODO reap old map entries, possibly eliminate the slice and return map values

	return bundles, nil
}
