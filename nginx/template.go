package nginx

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/30x/keymaster/client"
)

func Template(deploymentDir string, deployment *client.Deployment) *client.DeploymentError {

	pipesFileMap := make(map[string]string)                    // fq pipe name -> pipe file path
	bundleConfMap := make(map[*client.DeploymentBundle]string) // bundle -> bundle.conf path

	for _, b := range deployment.Bundles {
		bundlePath := path.Join(deploymentDir, b.BundleID)
		bundleConf := path.Join(bundlePath, "bundle.conf")
		bundle := b
		bundleConfMap[bundle] = bundleConf

		pipesDir := path.Join(bundlePath, "pipes")
		fis, err := ioutil.ReadDir(pipesDir)
		if err != nil {
			return &client.DeploymentError{ErrorCode: client.ErrorCodeTODO, Reason: err.Error()}
		}

		for _, fi := range fis {
			if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".yaml") {
				pipeName := strings.TrimSuffix(fi.Name(), ".yaml")
				fqPipeName := fmt.Sprintf("%s_%s", b.BundleID, pipeName)
				pipesFileMap[fqPipeName] = path.Join(pipesDir, fi.Name())
			}
		}
	}

	context := &TemplateContext{
		deploymentDir: deploymentDir,
		deployment:    deployment,
		bundleConfMap: bundleConfMap,
		pipesFileMap:  pipesFileMap, // fq pipe name -> pipe file path
	}

	err := templateNginxConf(context)
	if err != nil {
		return err
	}

	return templateBundleConfs(context)
}

func templateNginxConf(context *TemplateContext) *client.DeploymentError {

	nginxConf := path.Join(context.deploymentDir, "nginx.conf")
	return runTemplate(nginxConf, context)
}

func templateBundleConfs(context *TemplateContext) *client.DeploymentError {

	for _, bundle := range context.deployment.Bundles {
		bundleConf := context.bundleConfMap[bundle]
		bundleContext := DeploymentBundleContext{bundle: bundle, context: context}
		err := runTemplate(bundleConf, bundleContext)
		if err != nil {
			return err
		}
	}
	return nil
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

	// for debugging, writes file to stdout...
	//cmd := exec.Command("cat", file)
	//cmd.Stdout = os.Stdout
	//cmd.Run()

	return nil
}

type TemplateContext struct {
	deploymentDir string
	deployment    *client.Deployment
	bundleConfMap map[*client.DeploymentBundle]string
	bundleConfs   []string
	pipesFileMap  map[string]string // fq pipe name -> definition path
}

func (s TemplateContext) Bundles() (string, error) {

	tmpl := template.New("bundles")
	parsedTemplate, err := tmpl.Parse(bundlesTemplate)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = parsedTemplate.Execute(buf, s)

	return buf.String(), err
}

func (s TemplateContext) BundleNames() []string {
	return s.bundleConfs
}

func (s TemplateContext) Pipes() map[string]string {
	return s.pipesFileMap
}

type DeploymentBundleContext struct {
	bundle  *client.DeploymentBundle
	context *TemplateContext
}

func (c DeploymentBundleContext) Pipe(pipeName string) (string, error) {
	fqPipeName := fqPipeName(*c.bundle, pipeName)
	if c.context.pipesFileMap[fqPipeName] == "" { // no matching pipe
		return "", fmt.Errorf("No pipe named '%s' in bundle '%s'", pipeName, c.bundle.BundleID)
	}
	return fmt.Sprintf(pipeInclude, fqPipeName), nil
}

func fqPipeName(bundle client.DeploymentBundle, pipeName string) string {
	return fmt.Sprintf("%s_%s", bundle.BundleID, pipeName)
}

var pipeInclude = `
      set $goz_pipe '%s';
      include goz_pipe.conf;
`

var bundlesTemplate = `
  lua_package_path "./lua/?.lua;;;";

  init_worker_by_lua_block {
    libgozerian = require('lua-gozerian')
      -- pipes are assigned to yaml config urls (may be file: or http: refs)
      local pipes = {
      {{range $name, $file := .Pipes}}
        {{ $name }} = '{{ $file }}',
      {{end}}
      }

      libgozerian.init(pipes)
  }

  {{range $name := .BundleNames}}
    include {{ $name }}
  {{end}}
`
