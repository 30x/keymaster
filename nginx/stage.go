package nginx

import (
	"os"
	"path"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/util"
)

//StageManager the manager for staging a deployment
type StageManager interface {
	Stage(deployment *client.Deployment) (deploymentDir string, err *client.DeploymentError)
}

//StageManagerImpl impl placeholder impelemntation
type StageManagerImpl struct {
}

//Stage stage the manager
func (stageManager *StageManagerImpl) Stage(deployment *client.Deployment) (deploymentDir string, err *client.DeploymentError) {
	return Stage(deployment)
}

// Stage unzip, process templates, and validate the deployment.
// returns directory, DeploymentError
// if directory returned is not empty (may be non-empty even if error), client is responsible for cleanup
func Stage(deployment *client.Deployment) (string, *client.DeploymentError) {

	deploymentDir, err := util.MkTempDir("", deployment.ID, 0755)
	if err != nil {
		return "", &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	deploymentError := unzipSystem(deploymentDir, deployment)
	if deploymentError != nil {
		return deploymentDir, deploymentError
	}

	deploymentError = unzipDeploymentBundles(deploymentDir, deployment)
	if deploymentError != nil {
		return deploymentDir, deploymentError
	}

	// todo
	//deploymentError = Template(deploymentDir)
	//if deploymentError != nil {
	//	return deploymentDir, deploymentError
	//}

	// todo
	//deploymentError = ValidateDeployment(deploymentDir, deployment)
	//if deploymentError != nil {
	//	return deploymentDir, deploymentError
	//}

	return deploymentDir, deploymentError
}

// todo: may want to reconsider putting system at top level - possible name conflicts w/ deployment bundles?
func unzipSystem(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	err := util.Unzip(deployment.System.FilePath(), deploymentDir)
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	return nil
}

// unzipBundles unzip the deployment and return the directory
func unzipDeploymentBundles(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	for _, bundle := range deployment.Bundles {

		bundleDir := path.Join(deploymentDir, bundle.BundleID)
		err := os.Mkdir(bundleDir, 0755)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}

		err = util.Unzip(bundle.FilePath(), bundleDir)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}
	}

	return nil
}
