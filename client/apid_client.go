package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//ApidClient the client factory.  Use the CreateApidClient function to perform validation.
type ApidClient struct {
	apidHostPath      string
	downloadDirectory string
	client            *http.Client
}

//CreateApidClient create the client and validate the input
func CreateApidClient(apidHostPath string, downloadDirectory string) (*ApidClient, error) {
	dirInfo, err := os.Stat(downloadDirectory)

	if err != nil {

		//if it's a not exist error, we want to create it.  Otherwise, just return it
		if !os.IsNotExist(err) {
			return nil, err
		}

		err := os.MkdirAll(downloadDirectory, 0700)

		if err != nil {
			return nil, err
		}

	} else if !dirInfo.IsDir() {
		return nil, errors.New("Expected download directory to be a directory, it is not")
	}

	//return the apid client
	return &ApidClient{
		apidHostPath:      apidHostPath,
		downloadDirectory: downloadDirectory,
		client:            &http.Client{},
	}, nil
}

//DeploymentResponse the response to a deployment
type DeploymentResponse struct {
	DeploymentID string                     `json:"deploymentId"`
	Bundles      []DeploymentBundleResponse `json:"bundles"`
	ETag         string
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
	File *os.File
}

//PollDeployments poll the deployments fromthe apidHostPath with the etag (optional) and timeout (0 for none)
//returns the deployment response, or an error if one occurs.  A nil deploymentresponse indicates a timeout on polling (TODO, should this be a custom error?)
func (apidClient *ApidClient) PollDeployments(etag string, timeout int) (*DeploymentResponse, error) {

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

	deploymentResponse := &DeploymentResponse{
		ETag: responseETag,
	}

	err = json.NewDecoder(resp.Body).Decode(deploymentResponse)

	if err != nil {
		return nil, err
	}

	return deploymentResponse, nil
}

//GetBundle Get the bundle url result and write it to disk.  Returns a file pointer to the written file, or an error if the download did not occur.
func (apidClient *ApidClient) GetBundle(bundle DeploymentBundleResponse) (*DeploymentBundle, error) {
	req, err := http.NewRequest("GET", bundle.URL, nil)

	req.Header.Add("Accept", "application/zip")

	resp, err := apidClient.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Received response code %d with body %s", resp.StatusCode, body)
	}

	//we're good, copy the data into a file

	fileDir := filepath.FromSlash(apidClient.downloadDirectory +
		"/" + bundle.BundleID)

	file, err := os.Create(fileDir)

	if err != nil {
		return nil, err
	}

	io.Copy(file, resp.Body)

	deploymentBundle :=
		&DeploymentBundle{
			DeploymentBundleResponse: bundle,
			File: file,
		}

	return deploymentBundle, nil

}
