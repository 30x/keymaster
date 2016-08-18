package nginx

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//TestConfig Test the configuration of the nginx file.  Will return an error if an error or warning is detected
func TestConfig(prefixPath, configFile string) error {
	cmd := exec.Command("nginx", "-t", "-p", prefixPath, "-c", configFile)

	// log.Printf("About to execute command %+v", cmd)

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

//Start Start the nginx process with the prefix path, the config file path, and the start timeout.  If the start timeout elapses, a timeoutError will be thrown
func Start(prefixPath, configFilePath string, startTimeout time.Duration) error {

	command := exec.Command("nginx", "-p", prefixPath, "-c", configFilePath)

	log.Printf("About to start nginx with command %+v", command)

	//We have to grab our pipes beforehand, otherwise we wont' have a handle after start
	// stdErr, err := command.StderrPipe()

	stdErr := &bytes.Buffer{}
	stdOut := &bytes.Buffer{}

	command.Stderr = stdErr
	command.Stdout = stdOut

	//we have to use a function because command.Combine
	//start the command in a fork
	err := command.Start()

	if err != nil {
		log.Printf("Command was unable to start begin")
		defer log.Printf("Command was unable to start complete")
		return err
	}

	//Set a timer to timeout
	timer := time.AfterFunc(startTimeout, func() {
		log.Printf("Timeout occured when waiting for nginx start to exist after %s.  Pid is %d", startTimeout, command.ProcessState.Pid())
		// command.Process.Kill()
		// command.Process.Signal(syscall.SIGKILL)

		//TODO try to kill the nginx process from pid here, see what happens
	})

	//after error fails we want to stop the timer

	err = command.Wait()
	//stop the timer if the process has already finished
	timer.Stop()

	//now read stdout and stderr.  We deliberately discard errors.  In the event the process fails to start, these won't always exist
	stdOutString := string(stdOut.Bytes())
	stdErrString := string(stdErr.Bytes())

	log.Printf("Stdout :%s", stdOutString)

	log.Printf("Stderr :%s", stdErrString)

	if !command.ProcessState.Success() {
		err = errors.New("Process exiting with non 0 error code")
	}

	if err != nil {
		log.Printf("An error occured waiting for nginx start to complete.  Not sure it started. Error is %s", err)

		return &StartError{
			StdOut: stdOutString,
			StdErr: stdErrString,
			Err:    err,
		}
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

	pid, err := getNginxPid(pidFile)

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

//killProcess kill the process in the pid file if it exists
func killProcess(pidFile string) error {
	pid, err := getNginxPid(pidFile)

	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	//now send it signal 0 and see what happens

	err = process.Signal(syscall.SIGKILL)

	if err != nil {
		return err
	}

	return nil
}

func getNginxPid(pidFile string) (int, error) {
	_, err := os.Stat(pidFile)

	if err != nil {
		//if it's a not exist error, we swallow it, since it won't be running
		if err != os.ErrNotExist {
			return 0, err
		}

		return 0, nil
	}

	fileData, err := ioutil.ReadFile(pidFile)

	if err != nil {
		return 0, err
	}

	//nothing in the file, it's not running
	fileString := strings.TrimSpace(string(fileData))

	if len(fileString) == 0 {
		return 0, nil
	}

	pid, err := strconv.Atoi(fileString)

	if err != nil {
		return 0, err
	}

	return pid, err
}

//StartError an error where we failed to start nginx
type StartError struct {
	StdOut string
	StdErr string
	Err    error
}

func (e *StartError) Error() string {
	return fmt.Sprintf("Stdout is: \n%s\n\n Stdout is : \n%s\n\n.  Source error is %s", e.StdOut, e.StdErr, e.Err)
}
