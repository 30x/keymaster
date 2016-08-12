package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/nginx"
	"path"
	"io/ioutil"
	"github.com/30x/keymaster/util"
	"github.com/30x/keymaster/client"
)

var _ = Describe("stage", func() {

	Describe("unzip", func() {

		FIt("should unzip the bundles", func() {

			systemBundle := client.DeploymentBundle{
				BundleID: "bundle1",
				LocalFile: "../test/testsystem.zip",
			}

			bundles := make([]*client.DeploymentBundle, 1)
			bundles[0] = client.DeploymentBundle{
				BundleID: "bundle1",
				LocalFile: "../test/testbundle.zip",
			}

			deployment := client.Deployment{
				ID: "deployment_id",
				Bundles: bundles,
			}

			stageDir, err := nginx.Stage(deployment)
			Expect(err).NotTo(HaveOccurred())

			Expect(stageDir).Should(BeAnExistingFile())

			nginxConf:= path.Join(stageDir, "nginx.conf")
			Expect(nginxConf).Should(BeAnExistingFile())

			bundleDir := path.Join(stageDir, "bundle1")
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
