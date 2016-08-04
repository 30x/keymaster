package test

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/tnine/mockhttpserver"
)

const host = "localhost:9000"

//MockApidServer the container type for our test server
type MockApidServer struct {
	server *mockhttpserver.MockServer
}

//GetBundlesResponse the response json to the get bundles
type GetBundlesResponse struct {
	Bundles []Bundle `json:"bundles"`
}

//Bundle metadata information
type Bundle struct {
	BundleID string `json:"bundleId"`
	URL      string `json:"url"`
	AuthCode string `json:"authCode"`
}

//CreateMockApidServer create a mock apid server
func CreateMockApidServer() *MockApidServer {
	return &MockApidServer{
		server: &mockhttpserver.MockServer{},
	}
}

//CreateGetBundles Create a get bundle request that returns the specified http status and body.  Does not make use of the If-Non-Match or block headers.
func (mockServer *MockApidServer) CreateGetBundles(status int, bundles []Bundle, timeout int) error {

	response := GetBundlesResponse{
		Bundles: bundles,
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

//CreateGetBundle return the bundle bytes
func (mockServer *MockApidServer) CreateGetBundle(status int, bundleURL string, body []byte) {

	withoutServer := strings.Replace(bundleURL, "http://"+host, "", -1)
	mockServer.server.NewGet(withoutServer).ToResponse(status, body).Add()

}

//Start start the mock server
func (mockServer *MockApidServer) Start() {
	mockServer.server.StartAsync(host)
}

//Stop start the mock server
func (mockServer *MockApidServer) Stop() error {
	return mockServer.server.Shutdown()
}
