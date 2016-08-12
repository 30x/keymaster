package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/util"
	"io/ioutil"
	"path"
	"os"
)

var _ = Describe("Unzip", func() {

	It("should properly unzip a bundle file", func() {
		zipfile := "../test/testbundle.zip"

		tmpDir, err := ioutil.TempDir("", "TestUnzip")
		Expect(err).NotTo(HaveOccurred())

		defer os.RemoveAll(tmpDir)

		err = util.Unzip(zipfile, tmpDir)
		Expect(err).NotTo(HaveOccurred())

		bundleFile := path.Join(tmpDir, "bundle.conf")
		Expect(bundleFile).Should(BeAnExistingFile())

		pipeFile := path.Join(tmpDir, "pipes/apikey.yaml")
		Expect(pipeFile).Should(BeAnExistingFile())

		pipeFile = path.Join(tmpDir, "pipes/dump.yaml")
		Expect(pipeFile).Should(BeAnExistingFile())
	})
})
