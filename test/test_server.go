package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tnine/mockhttpserver"
)

const host = "localhost:9000"

//MockApidServer the container type for our test server
type MockApidServer struct {
	server *mockhttpserver.MockServer
}

//GetBundlesResponse the response json to the get bundles
type GetBundlesResponse struct {
	DeploymentID string       `json:"deploymentId"`
	System       SystemBundle `json:"system"`
	Bundles      []Bundle     `json:"bundles"`
}

//Bundle metadata information
type Bundle struct {
	SystemBundle
	AuthCode string `json:"authCode"`
}

//SystemBundle the system bundle
type SystemBundle struct {
	BundleID string `json:"bundleId"`
	URL      string `json:"url"`
}

//CreateMockApidServer create a mock apid server
func CreateMockApidServer() *MockApidServer {
	return &MockApidServer{
		server: &mockhttpserver.MockServer{},
	}
}

//CreateGetBundles Create a get bundle request that returns the specified http status and body.  Does not make use of the If-Non-Match or block headers.
func (mockServer *MockApidServer) CreateGetBundles(status int, deploymentID string, system SystemBundle, bundles []Bundle, timeout int) error {

	response := GetBundlesResponse{
		DeploymentID: deploymentID,
		Bundles:      bundles,
		System:       system,
	}

	data, err := json.Marshal(response)

	if err != nil {
		return err
	}

	log.Printf("Resposne data is %s", string(data))

	//set this up
	mockServer.server.NewGet("/deployments/current").ToResponse(status, data).Add()

	// if timeout > 0 {
	// 	timeoutValue := fmt.Sprintf("%d", timeout)
	// 	get = get.AddHeader("block", timeoutValue)
	// }

	return nil
}

//MockDeployment mock a response to the deployment
func (mockServer *MockApidServer) MockDeployment(deploymentID string, status int, body []byte) {
	url := fmt.Sprintf("/deployments/%s", deploymentID)
	mockServer.server.NewPost(url, "application/json", nil).ToResponse(status, body).Add()
}

//Start start the mock server
func (mockServer *MockApidServer) Start() {
	mockServer.server.StartAsync(host)
}

//Stop start the mock server
func (mockServer *MockApidServer) Stop() error {
	return mockServer.server.Shutdown()
}
