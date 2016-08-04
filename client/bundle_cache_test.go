package client_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle Tests", func() {

	It("Initial Load Empty Cache", func() {

		//start the mock api server

		timeout := 10

		bundleNames := []string{"1", "2", "3", "4", "5"}
		mockApiServer := createBundles(bundleNames, timeout)
		mockApiServer.Start()
		defer mockApiServer.Stop()

		//now load the cache, and ensure that we have everythig

		cache, err := client.CreateBundleCache("http://localhost:9000", "/tmp/initialloadbundles", timeout)

		Expect(err).Should(BeNil())

		bundles, err := cache.GetBundles()

		Expect(err).Should(BeNil())

		Expect(len(bundles)).Should(Equal(len(bundleNames)))

		for index, bundle := range bundles {

			// fmt.Printf("Bundle at index %d is +%v", index, *bundle)

			Expect(bundle).ShouldNot(BeNil())
			bundleID := fmt.Sprintf("%d", (index + 1))
			Expect(bundle.BundleID).Should(Equal(bundleID))
			Expect(len(bundle.AuthCode) > 0).Should(BeTrue())

			expected := fmt.Sprintf("http://localhost:9000/bundles/%s", bundleID)

			Expect(bundle.URL).Should(Equal(expected))

			Expect(bundle.File).ShouldNot(BeNil())

			file := bundle.File

			expectedName := fmt.Sprintf("/tmp/initialloadbundles/%s", bundleID)
			Expect(file.Name()).Should(Equal(expectedName))

		}

	})

	PIt("Load of completely new bundles", func() {

	})

	PIt("Partially new bundles", func() {

	})
})

func createBundles(bundleIds []string, timeout int) *test.MockApidServer {
	//pre allocate
	bundles := make([]test.Bundle, len(bundleIds))

	for i, bundleId := range bundleIds {
		url := fmt.Sprintf("http://localhost:9000/bundles/%s", bundleId)

		bundles[i] = test.Bundle{BundleID: bundleId, URL: url, AuthCode: fmt.Sprintf("%d", i)}
	}

	mockApiServer := test.CreateMockApidServer()

	mockApiServer.CreateGetBundles(http.StatusOK, bundles, timeout)

	//return the zip file in the test dir.  Eventually we'll want to change this
	for _, bundle := range bundles {
		bytes, err := ioutil.ReadFile("../test/testbundle.zip")

		Expect(err).Should(BeNil())

		//set up the bundle
		mockApiServer.CreateGetBundle(http.StatusOK, bundle.URL, bytes)
	}

	return mockApiServer

}
