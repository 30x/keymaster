package nginx_test

import (
	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/nginx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const nginxDir = ""
const basePath = "/Users/apigee/develop/go/src/github.com/30x/keymaster"

var _ = Describe("Manager", func() {

	FIt("Valid Configuration", func() {

		stager := &stageTester{
			testConfigDir: basePath + "/test/testbundles/validBundle",
		}

		// systemBundle := &client.SystemBundle{
		// 	BundleID: "bundle1",
		// 	URL:      "file://../test/testsystem.zip",
		// }

		// bundles := make([]*client.DeploymentBundle, 1)
		// bundles[0] = &client.DeploymentBundle{
		// 	BundleID: "bundle1",
		// 	URL:      "file://../test/testbundle.zip",
		// }

		deployment := &client.Deployment{
			ID: "deployment_id",
			// System:  systemBundle,
			// Bundles: bundles,
		}

		//wire up the resposne
		apiClient := &apiClientTester{
			mockDeployment: deployment,
		}

		manager := nginx.NewManager(apiClient, stager, nginxDir, 1)

		err := manager.ApplyDeployment()

		//no error with valid bundle
		Expect(err).Should(BeNil())

		//validate we returned successfully

		Expect(apiClient.deploymentResult.Status).Should(Equal(client.StatusSuccess))
		Expect(apiClient.deploymentResult.ID).Should(Equal(deployment.ID))

	})

})

//mock tester
type stageTester struct {
	//the dir to config, if no error is set, this is returned
	testConfigDir string
	//The error to set. If set it's returned.
	err *client.DeploymentError
}

func (test *stageTester) Stage(deployment *client.Deployment) (deploymentDir string, err *client.DeploymentError) {
	return test.testConfigDir, test.err
}

//mock tester
type apiClientTester struct {
	mockDeployment *client.Deployment

	pollDeploymentsErr error

	deploymentResult *client.DeploymentResult

	deploymentResultErr error
}

//PollDeployments poll the deployments and return the deployment
func (apiClient *apiClientTester) PollDeployments(etag string, timeout int) (*client.Deployment, error) {
	return apiClient.mockDeployment, apiClient.pollDeploymentsErr
}

//SetDeploymentResult set the deployment result
func (apiClient *apiClientTester) SetDeploymentResult(result *client.DeploymentResult) error {
	apiClient.deploymentResult = result
	return apiClient.deploymentResultErr
}
