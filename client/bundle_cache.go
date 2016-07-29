package keymaster

//BundleCache the type that represents the cache of bundles
type BundleCache struct {
	ApidURI string
	cache   []Bundle
}

//Bundle the type of the bundle to return
//TODO, not sure if this is even neccessary
type Bundle struct {
	DeploymentBundle
}

//GetBundles returns all bundles via a channel.  Blocks until a change in the bundle set is detected
func (bundleCache *BundleCache) GetBundles() ([]Bundle, error) {
	return []Bundle{}, nil
}
