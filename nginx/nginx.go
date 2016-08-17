package nginx

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

//TestConfig Test the configuration of the nginx file.  Will return an error if an error or warning is detected
func TestConfig(prefixPath, configFile string) error {
	cmd := exec.Command("nginx", "-t", "-p", prefixPath, "-c", configFile)

	log.Printf("About to execute command %+v", cmd)

	out, execErr := cmd.CombinedOutput()

	err := findError("emerg", out)
	if err != nil {
		return err
	}

	err = findError("warn", out)
	if err != nil {
		return err
	}

	//defer checking the execErr until the end.  Otherwise we won't receive emerg errors
	if execErr != nil {
		return execErr
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

//IsRunning return true if nginx is running, false otherwise.  A missing pid file is not considered an error
func IsRunning(pidFile string) (bool, error) {
	_, err := os.Stat(pidFile)

	if err != nil {
		//if it's a not exist error, we swallow it, since it won't be running
		if err != os.ErrNotExist {
			return false, err
		}

		return false, nil
	}

	fileData, err := ioutil.ReadFile(pidFile)

	if err != nil {
		return false, err
	}

	//nothing in the file, it's not running
	fileString := strings.TrimSpace(string(fileData))

	if len(fileString) == 0 {
		return false, nil
	}

	pid, err := strconv.Atoi(fileString)

	if err != nil {
		return false, err
	}

	process, err := os.FindProcess(pid)

	if err != nil {
		return false, err
	}

	//now send it signal 0 and see what happens

	err = process.Signal(syscall.Signal(0))

	if err != nil {
		return false, err
	}

	return true, nil

}
