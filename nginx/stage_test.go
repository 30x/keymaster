package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/nginx"
	"path"
	"github.com/30x/keymaster/client"
	"os"
)

var _ = Describe("stage", func() {

	Describe("unzip", func() {

		It("should unzip the bundles", func() {

			systemBundle := &client.SystemBundle{
				BundleID: "bundle1",
				URL: "file://../test/testsystem.zip",
			}

			bundles := make([]*client.DeploymentBundle, 1)
			bundles[0] = &client.DeploymentBundle{
				BundleID: "bundle1",
				URL: "file://../test/testbundle.zip",
			}

			deployment := &client.Deployment{
				ID: "deployment_id",
				System: systemBundle,
				Bundles: bundles,
			}

			stageDir, err := nginx.Stage(deployment)
			Expect(err).To(BeNil())
			Expect(stageDir).Should(BeAnExistingFile())

			defer os.RemoveAll(stageDir)

			nginxConf:= path.Join(stageDir, "nginx.conf")
			Expect(nginxConf).Should(BeAnExistingFile())

			bundleDir := path.Join(stageDir, bundles[0].BundleID)

			Expect(bundleDir).Should(BeAnExistingFile())

			nginxBundle := path.Join(bundleDir, "bundle.conf")
			Expect(nginxBundle).Should(BeAnExistingFile())

			pipeFile := path.Join(bundleDir, "pipes", "apikey.yaml")
			Expect(pipeFile).Should(BeAnExistingFile())

			pipeFile = path.Join(bundleDir, "pipes", "dump.yaml")
			Expect(pipeFile).Should(BeAnExistingFile())
		})
	})

})
