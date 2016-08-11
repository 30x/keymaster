package client_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle Tests", func() {

	It("Success", func() {

		deploymentID := "deploymentFoo"

		mockApiServer := test.CreateMockApidServer()
		mockApiServer.MockDeployment(deploymentID, http.StatusOK, nil)

		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now test it

		apiClient, err := client.CreateApidClient("http://localhost:9000")

		Expect(err).Should(BeNil())

		//
		deploymentResult := &client.DeploymentResult{
			ID:     deploymentID,
			Status: client.StatusSuccess,
		}

		err = apiClient.SetDeploymentResult(deploymentResult)

		Expect(err).Should(BeNil())
	})

	It("Parse Deployments", func() {

		//start the mock api server

		timeout := 10

		bundlePairs := []string{"1", "dependent", "2", "dependent", "3", "system", "4", "dependent", "5", "dependent"}

		deploymentID := "deployment1"

		mockApiServer := createBundles(deploymentID, bundlePairs, timeout)
		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now load the cache, and ensure that we have everythig

		client, err := client.CreateApidClient("http://localhost:9000")

		Expect(err).Should(BeNil())

		deployment, err := client.PollDeployments("", 60)

		Expect(err).Should(BeNil())

		Expect(err).Should(BeNil())

		Expect(deployment.ID).Should(Equal(deploymentID))

		bundles := deployment.Bundles

		Expect(len(bundles)).Should(Equal(len(bundlePairs) / 2))

		for index, bundle := range bundles {

			// fmt.Printf("Bundle at index %d is +%v", index, *bundle)

			Expect(bundle).ShouldNot(BeNil())
			bundleID := fmt.Sprintf("%d", (index + 1))
			Expect(bundle.BundleID).Should(Equal(bundleID))
			Expect(len(bundle.AuthCode) > 0).Should(BeTrue())

			expected := fmt.Sprintf("file:///tmp/keymaster-test/bundles/%s", bundleID)

			Expect(bundle.URL).Should(Equal(expected))

			Expect(bundle.LocalFile).Should(BeARegularFile())

			expectedName := fmt.Sprintf("/tmp/keymaster-test/bundles/%s", bundleID)
			Expect(bundle.LocalFile).Should(Equal(expectedName))

		}

	})

})

func createBundles(deploymentId string, bundlePairs []string, timeout int) *test.MockApidServer {
	//pre allocate
	bundles := []test.Bundle{}

	Expect(len(bundlePairs)%2).Should(Equal(0), "Pairs must be an even number")

	for i := 0; i < len(bundlePairs); i = i + 2 {

		bundleId := bundlePairs[i]
		bundleType := bundlePairs[i+1]

		url := fmt.Sprintf("file:///tmp/keymaster-test/bundles/%s", bundleId)

		testBundle := test.Bundle{BundleID: bundleId, URL: url, AuthCode: fmt.Sprintf("%d", i), Type: bundleType}
		bundles = append(bundles, testBundle)

		//copy the bundle over
		src, err := os.Open("../test/testbundle.zip")

		Expect(err).Should(BeNil())

		targetFile := strings.Replace(url, "file://", "", -1)

		targetDir := filepath.Dir(targetFile)

		err = os.MkdirAll(targetDir, 0770)

		Expect(err).Should(BeNil())

		target, err := os.Create(targetFile)

		Expect(err).Should(BeNil())

		_, err = io.Copy(src, target)

		Expect(err).Should(BeNil())
		src.Close()
		target.Close()

	}

	mockApiServer := test.CreateMockApidServer()

	mockApiServer.CreateGetBundles(http.StatusOK, deploymentId, bundles, timeout)

	return mockApiServer

}
