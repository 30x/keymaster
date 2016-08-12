package nginx

import (
	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/util"
	"os"
	"path"
)

type StageManager interface {
	Stage(deployment *client.Deployment) (deploymentDir string, err *client.DeploymentError)
}

// Stage unzip, process templates, and validate the deployment.
// returns directory, DeploymentError
func Stage(deployment *client.Deployment) (string, *client.DeploymentError) {

	deploymentDir, err := util.MkTempDir("", deployment.ID, 0755)
	if err != nil {
		return "", &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	deploymentError := processBundles(deploymentDir, deployment)
	if err != nil {
		return "", deploymentError
	}

	deploymentError = validateDeployment(deploymentDir, deployment)

	return deploymentDir, deploymentError
}


// unzipBundles unzip the deployment and return the directory
func processBundles(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	for _, bundle := range deployment.Bundles {

		bundleDir := path.Join(deploymentDir, bundle.BundleID)
		err := os.Mkdir(bundleDir, 0755)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}

		err = util.Unzip(bundle.LocalFile, bundleDir)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}

		templateBundle(bundleDir)

		//ValidateBundle(bundleDir) // todo
	}

	return nil
}

func templateBundle(bundleDir string) *client.DeploymentError {
	return nil
}

func validateDeployment(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	//bundleConfFile := path.Join(deploymentDir, "bundle.conf")
	//_, err := os.Stat(bundleConfFile)
	//if err != nil {
	//	return err
	//}

	// todo: validate pipes
	return nil
}
