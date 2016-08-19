package nginx_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/nginx"
	"github.com/30x/keymaster/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const nginxPidFile = "/usr/local/var/run/openresty.pid"
const basePath = "/Users/apigee/develop/go/src/github.com/30x/keymaster"

var _ = Describe("Manager", func() {

	// BeforeEach(func() {

	// 	nginxDir, err := util.MkTempDir("", "setup", 0755)

	// 	Expect(err).Should(BeNil())

	// 	defer os.RemoveAll(nginxDir)

	// 	//stop nginx, just in case it is running on the system
	// 	nginx.Stop(nginxDir)
	// })

	// AfterEach(func() {

	// 	nginxDir, err := util.MkTempDir("", "setup", 0755)

	// 	Expect(err).Should(BeNil())

	// 	defer os.RemoveAll(nginxDir)

	// 	//stop nginx, just in case it is running on the system
	// 	nginx.Stop(nginxDir)
	// })

	//tests a valid configuration on the first pass works
	FIt("Valid Configuration Single Pass", func() {

		fullPath, err := filepath.Abs("../test/testbundles/validBundle")

		Expect(err).Should(BeNil())

		stager := &stageTester{
			testConfigDir: fullPath,
		}

		deployment := &client.Deployment{
			ID: "deployment_id_1",
		}

		nginxDir, err := util.MkTempDir("", deployment.ID, 0755)

		Expect(err).Should(BeNil())

		defer nginx.Stop(nginxDir)
		defer os.RemoveAll(nginxDir)

		//wire up the resposne
		apiClient := &apiClientTester{
			mockDeployment: deployment,
		}

		manager := nginx.NewManager(apiClient, stager, nginxDir, nginxPidFile, 1)

		err = manager.ApplyDeployment()

		//no error with valid bundle
		Expect(err).Should(BeNil())

		//validate we returned successfully

		Expect(apiClient.deploymentResult.Status).Should(Equal(client.StatusSuccess))
		Expect(apiClient.deploymentResult.ID).Should(Equal(deployment.ID))

	})

	//tests a valid configuration on the first pass works
	It("Valid Configuration Multiple Pass", func() {

		nginxDir, err := util.MkTempDir("", "deployment_id_2", 0755)

		Expect(err).Should(BeNil())

		defer os.RemoveAll(nginxDir)

		//wire up the resposne
		apiClient := &apiClientTester{}

		stager := &stageTester{
			testConfigDir: basePath + "/test/testbundles/validBundle",
		}

		for i := 0; i < 5; i++ {

			//create a new deployment id each pass
			deploymentId := fmt.Sprintf("deployment_id_2_%d", i)
			deployment := &client.Deployment{
				ID: deploymentId,
			}

			apiClient.mockDeployment = deployment

			manager := nginx.NewManager(apiClient, stager, nginxDir, nginxPidFile, 1)

			err := manager.ApplyDeployment()

			//no error with valid bundle
			Expect(err).Should(BeNil())

			//validate we returned successfully

			Expect(apiClient.deploymentResult.Status).Should(Equal(client.StatusSuccess))
			Expect(apiClient.deploymentResult.ID).Should(Equal(deploymentId))
		}

	})

	//TODO, test success, fail, success

	It("Single Conflict Configuration", func() {

		stager := &stageTester{
			testConfigDir: basePath + "/test/testbundles/singleConflictPathBundle",
		}

		deployment := &client.Deployment{
			ID: "deployment_id_3",
		}

		nginxDir, err := util.MkTempDir("", deployment.ID, 0755)

		Expect(err).Should(BeNil())

		defer os.RemoveAll(nginxDir)

		//wire up the resposne
		apiClient := &apiClientTester{
			mockDeployment: deployment,
		}

		manager := nginx.NewManager(apiClient, stager, nginxDir, nginxPidFile, 1)

		err = manager.ApplyDeployment()

		//no error with valid bundle
		Expect(err).ShouldNot(BeNil())

		//check the error is applicable
		expectedErrorMessage := "[warn] conflicting server name \"localhost\" on 0.0.0.0:9000"
		containsMessage := strings.Contains(err.Error(), expectedErrorMessage)
		Expect(containsMessage).Should(BeTrue(), fmt.Sprintf("Should contain error message %s. Error message was %s", expectedErrorMessage, err))
		//validate we returned successfully

		Expect(apiClient.deploymentResult.Status).Should(Equal(client.StatusFail))
		Expect(apiClient.deploymentResult.ID).Should(Equal(deployment.ID))

	})

	//TODO, this isn't returning the error message, only an error code
	It("Multiple invalid files", func() {

		stager := &stageTester{
			testConfigDir: basePath + "/test/testbundles/multipleInvalidBundle",
		}

		deployment := &client.Deployment{
			ID: "deployment_id_4",
		}

		nginxDir, err := util.MkTempDir("", deployment.ID, 0755)

		Expect(err).Should(BeNil())

		defer os.RemoveAll(nginxDir)

		//wire up the resposne
		apiClient := &apiClientTester{
			mockDeployment: deployment,
		}

		manager := nginx.NewManager(apiClient, stager, nginxDir, nginxPidFile, 1)

		err = manager.ApplyDeployment()

		//no error with valid bundle
		Expect(err).ShouldNot(BeNil())

		//check the error is applicable
		expectedErrorMessage := "[emerg] directive \"listen\" is not terminated by \";\""
		log.Printf("Error is %s", err.Error())
		containsMessage := strings.Contains(err.Error(), expectedErrorMessage)
		Expect(containsMessage).Should(BeTrue(), fmt.Sprintf("Should contain error message %s. Error message was %s", expectedErrorMessage, err))
		//validate we returned successfully

		Expect(apiClient.deploymentResult.Status).Should(Equal(client.StatusFail))
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
