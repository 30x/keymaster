package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/30x/keymaster/nginx"
	"path"
	"github.com/30x/keymaster/client"
	"os"
	"io"
	"github.com/30x/keymaster/util"
)

var _ = Describe("templating", func() {

	It("should create a valid deployment", func() {

		systemBundle := &client.SystemBundle{
			BundleID: "bundle1",
			URL: "file://../test/testsystem.zip",
		}

		var bundles []*client.DeploymentBundle

		bundles = append(bundles, &client.DeploymentBundle{
			BundleID: "bundle1",
			URL: "file://../test/testbundle.zip",
		})

		bundles = append(bundles, &client.DeploymentBundle{
			BundleID: "bundle2",
			URL: "file://../test/testbundle.zip",
		})

		deployment := &client.Deployment{
			ID: "deployment_id",
			System: systemBundle,
			Bundles: bundles,
		}

		stageDir, err := util.MkTempDir("", deployment.ID, 0755)
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(stageDir)

		err = copyDirRecursive("../test/template/testsystem", stageDir)
		Expect(err).NotTo(HaveOccurred())

		nginxConf := path.Join(stageDir, "nginx.conf")
		Expect(nginxConf).Should(BeAnExistingFile())

		for _, b := range bundles {
			bundleDir := path.Join(stageDir, b.BundleID)
			err = copyDirRecursive("../test/template/testbundle", bundleDir)
			Expect(err).NotTo(HaveOccurred())
		}

		deploymentErr := nginx.Template(stageDir, deployment)
		Expect(deploymentErr).To(BeNil())

		err = nginx.TestConfig(stageDir, "nginx.conf")
		Expect(err).NotTo(HaveOccurred())
	})
})

func copyDirRecursive(source, dest string) error {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	fileInfos, err := directory.Readdir(-1)

	for _, fi := range fileInfos {

		sourceFilePath := source + "/" + fi.Name()
		destFilePath := dest + "/" + fi.Name()

		if fi.IsDir() {
			err = copyDirRecursive(sourceFilePath, destFilePath)
			if err != nil {
				break
			}
		} else {
			err = copyFile(sourceFilePath, destFilePath)
			if err != nil {
				break
			}
		}
	}

	return err
}

func copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		stat, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, stat.Mode())
		}

	}

	return err
}