package nginx

import "github.com/30x/keymaster/client"

//Manager The config manager
type Manager struct {
	bundleCache        *client.BundleCache
	lastApidDeployment *client.Deployment

	lastUnzippedDeployment *UnzippedDeployment
}

//UnzippedDeployment a struct representing a deployment that has been unzipped
type UnzippedDeployment struct {
	targetDir string
	//targetBundles some kind of object per bundle?

}

//NewManager Create a new instance of the configuration manager
func NewManager(bundleCache *client.BundleCache) *Manager {
	return &Manager{
		bundleCache: bundleCache,
	}
}

//ApplyDeployment Runs once, attempting to apply the latest deployment from the bundle cache. May return an execution error if there is a problem executing
func (manager *Manager) ApplyDeployment() error {

	deployment, err := manager.bundleCache.GetBundles()

	if err != nil {
		return err
	}

	//same deployment as last time, do nothing
	if manager.lastApidDeployment != nil && deployment.ID == manager.lastApidDeployment.ID {
		return nil
	}

	//we have a new deployment, time to apply it
	//

	//unzip bundles to bundle id directlry

	unzipped, unzipErr := UnzipBundle(deployment)

	if unzipErr != nil {
		return unzipErr
	}

	//perform template processing

	processingError := ProcessTemplates(unzipped)

	if processingError != nil {
		return processingError
	}

	//test nginx with the processed templates/new configs.  TODO warnings constitute a failure

	//reload nginx

	//reset pointers to last

	//cleanup old last from file system

	//TODO add a template where the deployment.ID is returned at localhost:5280/ to validate we're actually running and get the status of the system

	return nil

}

func (manager *Manager) deploymentFailed(deployment *client.Deployment) {

}

func (manager *Manager) deploymentComplete(deployment *client.Deployment) {

}

//UnzipBundle  unzip the deployment and return the struct with the info for the directory and bundle
func UnzipBundle(deployment *client.Deployment) (UnzippedDeployment, error) {
	return UnzippedDeployment{}, nil
}

//ProcessTemplates Process the bundle templates.  Return an error if one occurs
func ProcessTemplates(unzippedDeployment UnzippedDeployment) error {
	return nil
}
