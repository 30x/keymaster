package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//ApidClient the client factory.  Use the CreateApidClient function to perform validation.
type ApidClient struct {
	apidHostPath string
	client       *http.Client
}

//DeploymentType the type of demployment
type DeploymentType string

const (
	//SystemType the type of bundle for a system bundle
	SystemType = "system"

	//DependentType the type of the bundle when it is a user defined type
	DependentType = "dependent"
)

//Deployment the type of deployment to return
type Deployment struct {
	ETAG    string
	ID      string              `json:"deploymentId"`
	Bundles []*DeploymentBundle `json:"bundles"`
}

//DeploymentBundle the bundle to deploy in a response
type DeploymentBundle struct {
	BundleID string         `json:"bundleId"`
	AuthCode string         `json:"authCode"`
	URL      string         `json:"url"`
	Type     DeploymentType `json:"type"`

	//the path on the local system to the file in the url
	LocalFile string
}

//DeploymentResult the result of a deployment
type DeploymentResult struct {
	//The deploymentId
	ID string
	//The status of the deployment
	Status DeploymentStatus `json:"status"`
	//Any errors that may have occurred.  If we're successful, this can be nil or empty
	Error DeploymentError
}

//DeploymentError The error that occurred on deployment
type DeploymentError struct {
	ErrorCode string `json:"errorCode"`
	Reason    string `json:"reason"`
	bundleErrors []BundleError
}

//BundleError Any Bundle-specific error that occurred on deployment
type BundleError struct {
	BundleID string `json:"bundleId"`
	ErrorCode string `json:"errorCode"`
	Reason    string `json:"reason"`
}

//DeploymentStatus the status of the deployment
type DeploymentStatus string

const (
	//StatusFail the deployment failed
	StatusFail DeploymentStatus = "FAIL"
	//StatusSuccess the deployment succeeded.
	StatusSuccess DeploymentStatus = "SUCCESS"
)

//CreateApidClient create the client and validate the input
func CreateApidClient(apidHostPath string) (*ApidClient, error) {

	//return the apid client
	return &ApidClient{
		apidHostPath: apidHostPath,
		client:       &http.Client{},
	}, nil
}

//PollDeployments poll the deployments fromthe apidHostPath with the etag (optional) and timeout (0 for none)
//returns the deployment response, or an error if one occurs.  A nil deploymentresponse indicates a timeout on polling (TODO, should this be a custom error?)
func (apidClient *ApidClient) PollDeployments(etag string, timeout int) (*Deployment, error) {

	url := apidClient.apidHostPath + "/deployments/current"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if len(etag) > 0 {
		req.Header.Add("If-None-Match", etag)
	}

	if timeout > 0 {
		req.Header.Add("block", string(timeout))
	}
	req.Header.Add("Accept", "application/zip")

	resp, err := apidClient.client.Do(req)

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

	responseETag := resp.Header.Get("ETag")

	deploymentResponse := &Deployment{
		ETAG: responseETag,
	}

	err = json.NewDecoder(resp.Body).Decode(deploymentResponse)

	if err != nil {
		return nil, err
	}

	//link up the files
	for _, bundle := range deploymentResponse.Bundles {
		bundle.LocalFile = strings.Replace(bundle.URL, "file://", "", -1)
	}

	return deploymentResponse, nil
}

//SetDeploymentResult set the result of the deployment.  Returns an error if the call was unsuccessful
func (apidClient *ApidClient) SetDeploymentResult(result *DeploymentResult) error {

	url := fmt.Sprintf("%s/deployments/%s", apidClient.apidHostPath, result.ID)

	payload, err := json.Marshal(result)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

	req.Header.Add("Content-Type", "application/json")

	resp, err := apidClient.client.Do(req)

	if err != nil {
		return err
	}

	//if it wasn't successful, throw an error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected response code %d, but response code was %d.  Reason is %s", http.StatusOK, resp.StatusCode, resp.Status)
	}

	//no need to read the body
	resp.Body.Close()

	return nil
}
