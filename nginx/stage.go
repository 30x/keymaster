package nginx

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/util"
)

//Stage unzip, process templates, and validate the deployment
func Stage(deployment *client.Deployment) (extractedDirectory string, errors []*client.DeploymentError) {
	return "", nil
}

//UnzipBundle  unzip the deployment and return the struct with the info for the directory and bundle
func UnzipBundle(deployment *client.Deployment) (string, error) {

	deploymentDir, err := ioutil.TempDir("", "deployment_"+deployment.ID)
	if err != nil {
		return "", err
	}

	for _, bundle := range deployment.Bundles {

		bundleDir, err := ioutil.TempDir(deploymentDir, "bundle_"+bundle.BundleID)

		err = util.Unzip(bundle.LocalFile, bundleDir)
		if err != nil {
			return "", err // todo: specific err identifying bundle
		}

		ValidateBundle(bundleDir) // todo: specific err identifying bundle
	}

	return deploymentDir, nil
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
func ProcessTemplates(unzippedDir string) error {
	return nil
}
