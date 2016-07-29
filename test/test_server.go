package test

import (
	"encoding/json"

	"github.com/tnine/mockhttpserver"
)

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
}

//CreateMockApidServer create a mock apid server
func CreateMockApidServer() *MockApidServer {
	return &MockApidServer{
		server: &mockhttpserver.MockServer{},
	}
}

//CreateGetBundles Create a get bundle request that returns the specified http status and body.  Does not make use of the If-Non-Match or block headers.
func (mockServer *MockApidServer) CreateGetBundles(status int, bundles []Bundle) error {

	response := GetBundlesResponse{
		Bundles: bundles,
	}

	data, err := json.Marshal(response)

	if err != nil {
		return err
	}

	//set this up
	mockServer.server.NewGet("/bundles").ToResponse(status, data)

	return nil
}

//CreateGetBundle return the bundle bytes
func (mockServer *MockApidServer) CreateGetBundle(status int, bundleURL string, body []byte) {

	mockServer.server.NewGet(bundleURL).AddHeader("Accept", "application/x-tar").ToResponse(status, body)

}

//Start start the mock server
func (mockServer *MockApidServer) Start() {
	mockServer.server.Listen("localhost:")
}

//Stop start the mock server
func (mockServer *MockApidServer) Stop() error {
	return mockServer.server.Shutdown()
}
