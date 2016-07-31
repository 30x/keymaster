package nginx

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"fmt"
)

func TestConfig(configFile string) error {
	cmd := exec.Command("nginx", "-t", "-p", filepath.Dir(configFile), "-c", configFile)

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
	matched, err := regexp.Match("^nginx: \\[" + errType + "\\]", message)
	if err != nil {
		return err
	}
	if matched {
		return fmt.Errorf("Config error:\n%s", message)
	}
	return nil
}
