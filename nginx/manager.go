package nginx

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/30x/keymaster/client"
)

//Manager The config manager
type Manager struct {
	client       client.ApidClient
	nginxPidFile string
	nginxWorkDir string
	pollTimeout  int
	stageManager StageManager

	//state of last successful deployment
	lastApidDeployment     *client.Deployment
	lastUnzippedDeployment string
}

//NewManager Create a new instance of the configuration manager
func NewManager(apiClient client.ApidClient, stageManager StageManager, nginxWorkDir string, nginxPidFile string, pollTimeout int) *Manager {
	return &Manager{
		client:       apiClient,
		stageManager: stageManager,
		nginxWorkDir: nginxWorkDir,
		nginxPidFile: nginxPidFile,
		pollTimeout:  pollTimeout,
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

	unzippedDir, deploymentError := manager.stageManager.Stage(deployment)

	if deploymentError != nil {
		manager.setDeploymentStatus(deployment, deploymentError)
		return err
	}

	//perform template processing

	//test nginx with the processed templates/new configs.  TODO warnings constitute a failure

	systemFile := fmt.Sprintf("%s/%s", unzippedDir, "nginx.conf")
	err = TestConfig(manager.nginxWorkDir, systemFile)

	if err != nil {
		manager.signalError(deployment, err)
		return err
	}

	//reload or start nginx if not running
	//TODO detect start state from PID

	isRunning, err := IsRunning(manager.nginxPidFile)

	if err != nil {
		manager.signalError(deployment, err)

		return err
	}

	if isRunning {
		err = Reload(manager.nginxWorkDir, systemFile)
	} else {
		err = Start(manager.nginxWorkDir, systemFile, 5*time.Second)
	}

	if err != nil {
		manager.signalError(deployment, err)

		return err
	}

	//reset pointers to last for our next invocation

	previousUnzipped := manager.lastUnzippedDeployment

	manager.lastApidDeployment = deployment
	manager.lastUnzippedDeployment = unzippedDir

	//cleanup old last from file system
	err = os.RemoveAll(previousUnzipped)

	//swallow this error, it shouldn't blow up our process
	if err != nil {
		log.Printf("Unable to remove directory %s.  Error is %s", previousUnzipped, err)
	}

	manager.setDeploymentStatus(deployment, nil)
	//TODO add a template where the deployment.ID is returned at localhost:5280/ to validate we're actually running and get the status of the system

	return nil

}

func (manager *Manager) signalError(deployment *client.Deployment, err error) {
	deploymentError := &client.DeploymentError{
		ErrorCode: client.ErrorCodeTODO,
		Reason:    err.Error(),
	}

	manager.setDeploymentStatus(deployment, deploymentError)

}

func (manager *Manager) setDeploymentStatus(deployment *client.Deployment, err *client.DeploymentError) {
	deploymentResult := &client.DeploymentResult{
		ID:     deployment.ID,
		Status: client.StatusSuccess,
	}

	if err != nil {

		deploymentResult.Error = err
		deploymentResult.Status = client.StatusFail

	}

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
