package nginx

import (
	"fmt"

	"log"
	"os"

	"github.com/30x/keymaster/client"
)

//Manager The config manager
type Manager struct {
	client      *client.ApidClient
	nginxDir    string
	pollTimeout int

	//state of last successful deployment
	lastApidDeployment     *client.Deployment
	lastUnzippedDeployment *UnzippedDeployment
}

//UnzippedDeployment a struct representing a deployment that has been unzipped
type UnzippedDeployment struct {
	targetDir string
}

//NewManager Create a new instance of the configuration manager
func NewManager(apiClient *client.ApidClient, nginxDir string, pollTimeout int) *Manager {
	return &Manager{
		client:      apiClient,
		nginxDir:    nginxDir,
		pollTimeout: pollTimeout,
	}
}

//ApplyDeployment Runs once, attempting to apply the latest deployment from the bundle cache. May return an execution error if there is a problem executing
func (manager *Manager) ApplyDeployment() error {

	etag := ""

	if manager.lastApidDeployment != nil {
		etag = manager.lastApidDeployment.ETAG
	}

	deployment, err := manager.client.PollDeployments(etag, manager.pollTimeout)

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

	systemFile := fmt.Sprintf("%s/%s", unzipped.targetDir, "nginx.conf")
	err = TestConfig(systemFile)

	if err != nil {
		manager.deploymentFailed(deployment, []*client.DeploymentBundle{}, err)
		return err
	}

	//reload or start nginx if not running
	//TODO detect start state from PID

	err = Reload(manager.nginxDir, systemFile)

	if err != nil {
		manager.deploymentFailed(deployment, []*client.DeploymentBundle{}, err)
		return err
	}

	//reset pointers to last for our next invocation

	previousUnzipped := manager.lastUnzippedDeployment

	manager.lastApidDeployment = deployment
	manager.lastUnzippedDeployment = &unzipped

	//cleanup old last from file system

	err = os.RemoveAll(previousUnzipped.targetDir)

	//swallow this error, it shouldn't blow up our process
	if err != nil {
		log.Printf("Unable to remove directory %s.  Error is %s", previousUnzipped.targetDir, err)
	}

	//TODO add a template where the deployment.ID is returned at localhost:5280/ to validate we're actually running and get the status of the system

	return nil

}

func (manager *Manager) deploymentFailed(deployment *client.Deployment, failedBundles []*client.DeploymentBundle, err error) {
	deploymentResult := &client.DeploymentResult{
		ID: deployment.ID,
	}

	status := client.StatusSuccess
	errorMessage := ""

	if err != nil {
		errorMessage = err.Error()
		status = client.StatusFail
	}

	deploymentResult.Errors = []*client.DeploymentError{}

	for _, failedBundle := range failedBundles {
		deployError := &client.DeploymentError{
			//todo, how can we tell which bundle failed?
			BundleID:  failedBundle.BundleID,
			ErrorCode: errorMessage,
			Reason:    errorMessage,
		}

		deploymentResult.Errors = append(deploymentResult.Errors, deployError)
	}

	deploymentResult.Status = status

	setErr := manager.client.SetDeploymentResult(deploymentResult)

	if setErr != nil {
		log.Printf("Error calling apid. Not setting failure %s", setErr)
		//TODO if we can't set our status, should we fail here and restart?
	}
}

func (manager *Manager) deploymentComplete(deployment *client.Deployment, err error) {
	deploymentResult := &client.DeploymentResult{
		ID: deployment.ID,
	}

	setErr := manager.client.SetDeploymentResult(deploymentResult)

	if setErr != nil {
		log.Printf("Error calling apid. Not setting success %s", setErr)
		//TODO if we can't set our status, should we fail here and restart?
	}
}

//UnzipBundle  unzip the deployment and return the struct with the info for the directory and bundle
func UnzipBundle(deployment *client.Deployment) (UnzippedDeployment, error) {
	return UnzippedDeployment{}, nil
}

//ProcessTemplates Process the bundle templates.  Return an error if one occurs
func ProcessTemplates(unzippedDeployment UnzippedDeployment) error {
	return nil
}
