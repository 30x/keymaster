package nginx

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

//TestConfig Test the configuration of the nginx file.  Will return an error if an error or warning is detected
func TestConfig(prefixPath, configFile string) error {
	cmd := exec.Command("nginx", "-t", "-p", prefixPath, "-c", configFile)

	log.Printf("About to execute command %+v", cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	err = findError("emerg", out)
	if err != nil {
		return err
	}

	err = findError("warn", out)
	if err != nil {
		return err
	}

	return nil
}

func findError(errType string, message []byte) error {
	matched, err := regexp.Match("^nginx: \\["+errType+"\\]", message)
	if err != nil {
		return err
	}
	if matched {
		return fmt.Errorf("Config error:\n%s", message)
	}
	return nil
}

func Start(prefixPath, configFilePath string) error {
	out, err := exec.Command("nginx", "-p", prefixPath, "-c", configFilePath).CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}
	return nil
}

func Stop(prefixPath string) error {
	out, err := exec.Command("nginx", "-p", prefixPath, "-s", "stop").CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}
	return nil
}

func Reload(prefixPath, configFilePath string) error {
	out, err := exec.Command("nginx", "-p", prefixPath, "-c", configFilePath, "-s", "reload").CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}
	return nil
}
