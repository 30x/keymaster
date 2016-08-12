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

		bundleIds := []string{"1", "2", "3", "4", "5"}

		deploymentID := "deployment1"

		mockApiServer := createBundles(deploymentID, bundleIds, timeout)
		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now load the cache, and ensure that we have everythig

		client, err := client.CreateApidClient("http://localhost:9000")

		Expect(err).Should(BeNil())

		deployment, err := client.PollDeployments("", 60)

		Expect(err).Should(BeNil())

		Expect(deployment.ID).Should(Equal(deploymentID))

		bundles := deployment.Bundles

		Expect(len(bundles)).Should(Equal(len(bundleIds)))

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

		Expect(deployment.System).ShouldNot(BeNil())

		Expect(deployment.System.BundleID).Should(Equal("system-revision-1"))
		Expect(deployment.System.URL).Should(Equal("file:///tmp/keymaster-test/bundles/system-revision-1"))
		Expect(deployment.System.LocalFile).Should(BeARegularFile())

	})

})

func createBundles(deploymentId string, bundIds []string, timeout int) *test.MockApidServer {
	//pre allocate
	bundles := []test.Bundle{}

	for _, bundleId := range bundIds {

		url := fmt.Sprintf("file:///tmp/keymaster-test/bundles/%s", bundleId)

		testBundle := test.Bundle{}
		testBundle.BundleID = bundleId
		testBundle.URL = url
		testBundle.AuthCode = bundleId

		bundles = append(bundles, testBundle)

		copyBundleFile(url)

	}

	mockApiServer := test.CreateMockApidServer()

	systemBundle := test.SystemBundle{
		BundleID: "system-revision-1",
		URL:      "file:///tmp/keymaster-test/bundles/system-revision-1",
	}

	copyBundleFile(systemBundle.URL)

	mockApiServer.CreateGetBundles(http.StatusOK, deploymentId, systemBundle, bundles, timeout)

	return mockApiServer

}

func copyBundleFile(destFileUrl string) {
	//copy the bundle over
	src, err := os.Open("../test/testbundle.zip")

	Expect(err).Should(BeNil())

	targetFile := strings.Replace(destFileUrl, "file://", "", -1)

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
