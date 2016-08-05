package client_test

import (
	"net/http"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle Tests", func() {

	FIt("Success", func() {

		deploymentID := "deploymentFoo"

		mockApiServer := test.CreateMockApidServer()
		mockApiServer.MockDeployment(deploymentID, http.StatusOK, nil)

		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now test it

		apiClient, err := client.CreateApidClient("http://localhost:9000", "/tmp/mytmpdir")

		Expect(err).Should(BeNil())

		//
		deploymentResult := &client.DeploymentResult{
			ID:     deploymentID,
			Status: client.StatusSuccess,
		}

		err = apiClient.SetDeploymentResult(deploymentResult)

		Expect(err).Should(BeNil())
	})
})
