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

	FIt("Initial Load Empty Cache", func() {

		//start the mock api server

		timeout := 10

		bundleNames := []string{"1", "2", "3", "4", "5"}

		deploymentID := "deployment1"

		mockApiServer := createBundles(deploymentID, bundleNames, timeout)
		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now load the cache, and ensure that we have everythig

		cache, err := client.CreateBundleCache("http://localhost:9000", "/tmp/initialloadbundles", timeout)

		Expect(err).Should(BeNil())

		deployment, err := cache.GetBundles()

		Expect(err).Should(BeNil())

		Expect(deployment.ID).Should(Equal(deploymentID))

		bundles := deployment.Bundles

		Expect(len(bundles)).Should(Equal(len(bundleNames)))

		for index, bundle := range bundles {

			// fmt.Printf("Bundle at index %d is +%v", index, *bundle)

			Expect(bundle).ShouldNot(BeNil())
			bundleID := fmt.Sprintf("%d", (index + 1))
			Expect(bundle.BundleID).Should(Equal(bundleID))
			Expect(len(bundle.AuthCode) > 0).Should(BeTrue())

			expected := fmt.Sprintf("file:///tmp/keymaster-test/bundles/%s", bundleID)

			Expect(bundle.URL).Should(Equal(expected))

			Expect(bundle.File).ShouldNot(BeNil())

			file := bundle.File

			expectedName := fmt.Sprintf("/tmp/keymaster-test/bundles/%s", bundleID)
			Expect(file.Name()).Should(Equal(expectedName))

		}

	})

	It("Ack client bundle", func() {

	})

	PIt("Load of completely new bundles", func() {

	})

	PIt("Partially new bundles", func() {

	})
})

func createBundles(deploymentId string, bundleIds []string, timeout int) *test.MockApidServer {
	//pre allocate
	bundles := make([]test.Bundle, len(bundleIds))

	for i, bundleId := range bundleIds {
		url := fmt.Sprintf("file:///tmp/keymaster-test/bundles/%s", bundleId)

		bundles[i] = test.Bundle{BundleID: bundleId, URL: url, AuthCode: fmt.Sprintf("%d", i)}

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
