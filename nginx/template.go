package nginx

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/30x/keymaster/client"
	"gopkg.in/yaml.v2"
)

type bundleMetadataDef struct {
	Pipes map[string]string `json:"pipes"`
}

func Template(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	bundles := make(map[string]bundle)

	for _, b := range deployment.Bundles {
		bundlePath := path.Join(deploymentDir, b.BundleID)

		yamlBytes, err := ioutil.ReadFile(path.Join(bundlePath, "bundle.yaml"))
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}
		bundleMetadata := bundleMetadataDef{}
		err = yaml.Unmarshal(yamlBytes, &bundleMetadata)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}
		pipePaths := bundleMetadata.Pipes
		if pipePaths == nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}
		pathPipes := make(map[string]string)
		for k, v := range pipePaths {
			pathPipes[v] = k
		}

		pipesDir := path.Join(bundlePath, "pipes")
		fis, err := ioutil.ReadDir(pipesDir)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}

		pipes := make(map[string]pipe)
		for _, fileInfo := range fis {
			if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".yaml") {
				filePath := path.Join(pipesDir, fileInfo.Name())
				pipeName := strings.TrimSuffix(fileInfo.Name(), ".yaml")
				pipePath := pathPipes[pipeName]

				if pipePath == "" {
					errMsg := fmt.Sprintf("Pipe named %s does not exist in bundle %s", pipeName, b.BundleID)
					return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: errMsg}
				}

				p := pipe{
					filePath: filePath,
					Name:     pipeName,
					Path:     pipePath,
				}

				pipes[pipeName] = p
			}
		}

		bn := bundle{
			bundlePath:   bundlePath,
			VirtualHosts: b.VirtualHosts,
			Basepath:     b.BasePath,
			Target:       b.Target,
			Pipes:        pipes,
		}
		bundles[b.BundleID] = bn
	}

	nginxConfContext := &templateContext{
		deployment:    deployment,
		deploymentDir: deploymentDir,
		Bundles:       bundles,
	}

	nginxConfTemplate := path.Join(deploymentDir, "nginx.conf")
	return runTemplate(nginxConfTemplate, nginxConfContext)
}

func runTemplate(fileName string, context interface{}) *client.DeploymentError {

	parsedTemplate, err := template.ParseFiles(fileName)
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	writer := bufio.NewWriter(file)
	err = parsedTemplate.Execute(writer, context)
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}
	err = writer.Flush()
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}
	err = file.Close()
	if err != nil {
		return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
	}

	return nil
}

type pipe struct {
	filePath string

	Name string
	Path string
}

type bundle struct {
	bundlePath string

	VirtualHosts []string
	Basepath     string
	Target       string
	Pipes        map[string]pipe
}

type templateContext struct {
	deployment    *client.Deployment
	deploymentDir string

	Bundles map[string]bundle
}
