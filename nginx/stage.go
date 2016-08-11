package nginx

import (
	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/util"
	"io/ioutil"
	"os"
	"path"
)

//UnzipBundle  unzip the deployment and return the struct with the info for the directory and bundle
func UnzipBundle(deployment *client.Deployment) (*UnzippedDeployment, error) {

	deploymentDir, err := ioutil.TempDir("", "deployment_"+deployment.ID)
	if err != nil {
		return nil, err
	}

	for _, bundle := range deployment.Bundles {

		bundleDir, err := ioutil.TempDir(deploymentDir, "bundle_"+bundle.BundleID)

		err = util.Unzip(bundle.File.Name(), bundleDir)
		if err != nil {
			return nil, err // todo: specific err identifying bundle
		}

		ValidateBundle(bundleDir) // todo: specific err identifying bundle
	}

	return &UnzippedDeployment{targetDir: deploymentDir}, nil
}

func ValidateBundle(bundleDir string) error {

	bundleConfFile := path.Join(bundleDir, "bundle.conf")
	_, err := os.Stat(bundleConfFile)
	if err != nil {
		return err
	}

	// todo: validate pipes
	return nil
}

//ProcessTemplates Process the bundle templates.  Return an error if one occurs
func ProcessTemplates(unzippedDeployment *UnzippedDeployment) error {
	return nil
}
