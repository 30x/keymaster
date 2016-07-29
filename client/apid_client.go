package keymaster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//DeploymentResponse the response to a deployment
type DeploymentResponse struct {
	DeploymentID string                     `json:"deploymentId"`
	Bundles      []DeploymentBundleResponse `json:"bundles"`
}

//DeploymentBundleResponse the bundle to deploy in a response
type DeploymentBundleResponse struct {
	BundleID string `json:"bundleId"`
	AuthCode string `json:"authCode"`
	URL      string `json:"url"`
}

//DeploymentBundle The actual deployment bundle with the file data present
type DeploymentBundle struct {
	DeploymentBundleResponse
	file os.File
}

//PollDeployments poll the deployments fromthe apidHostPath with the etag (optional) and timeout (0 for none)
//returns the deployment response, or an error if one occurs.  A nil deploymentresponse indicates a timeout on polling (TODO, should this be a custom error?)
func PollDeployments(apidHostPath string, etag string, timeout int) (*DeploymentResponse, error) {

	client := &http.Client{}

	url := apidHostPath + "/deployments/current"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	//we timed out, return nothing.  TODO make this a better error type
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Could not poll deployments Status code is %d with body %s.", resp.StatusCode, string(errorBody))
	}

	deploymentResponse := &DeploymentResponse{}

	err = json.NewDecoder(resp.Body).Decode(deploymentResponse)

	if err != nil {
		return nil, err
	}

	return deploymentResponse, nil
}

//GetBundle Get the bundle url result and write it to disk.  Returns a file pointer to the written file, or an error if the download did not occur.
func GetBundle(bundle *DeploymentBundleResponse) (*DeploymentBundle, error) {
	return nil, nil

}
