package terraform

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/machine/libmachine/log"
)

// Terraformer is used to invoke Terraform.
type Terraformer struct {
	// The path to the Terraform executable
	ExecutablePath string

	// The directory containing the Terraform configuration
	ConfigDir string
}

// New creates a new Terraformer using the specified configuration directory.
func New(configDir string) (*Terraformer, error) {
	terraformer := &Terraformer{
		ConfigDir: configDir,
	}
	err := terraformer.resolveExecutablePath()
	if err != nil {
		return nil, err
	}

	return terraformer, nil
}

// NewWithExecutable creates a new Terraformer using the specified Terraform executable and configuration directory.
func NewWithExecutable(executablePath string, configDir string) (*Terraformer, error) {
	terraformer := &Terraformer{
		ExecutablePath: executablePath,
		ConfigDir:      configDir,
	}
	err := terraformer.resolveExecutablePath()
	if err != nil {
		return nil, err
	}

	return terraformer, nil
}

// Run invokes Terraform.
//
// command is the name of the Terraform command to execute (e.g. plan, apply, output,  destroy, etc)
// arguments are any other arguments to pass to Terraform
func (terraformer *Terraformer) Run(command string, arguments ...string) (success bool, output string, err error) {
	var executablePath string
	executablePath, err = terraformer.getExecutablePath()
	if err != nil {
		return
	}

	args := []string{command}
	args = append(args, arguments...)

	commandLine := fmt.Sprintf("%s %s",
		terraformer.ExecutablePath,
		command,
	)
	if len(arguments) > 0 {
		commandLine += " "
		commandLine += strings.Join(arguments, " ")
	}

	var (
		terraformCommand *exec.Cmd
		programOutput    bytes.Buffer
	)
	terraformCommand = exec.Command(executablePath, args...)
	terraformCommand.Stdout = &programOutput
	terraformCommand.Stderr = &programOutput
	terraformCommand.Dir = terraformer.ConfigDir // Always run Terraform in the configuration directory

	log.Debugf(`Executing "%s" ...`, commandLine)
	err = terraformCommand.Run()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Failed (%s)", err.Error())

		output = string(
			programOutput.Bytes(),
		)

		return
	}

	success = true

	output = string(
		programOutput.Bytes(),
	)

	log.Debugf("Successfully executed %s ...", commandLine)

	return
}

// PipeHandler is a function which receives lines of piped output from Terraform as they become available
type PipeHandler func(string)

// RunPiped invokes Terraform and pipes its output to the specified OutputHandler.
//
// command is the name of the Terraform command to execute (e.g. plan, apply, output,  destroy, etc)
// pipeOutput is a function called once for each line of output received
// arguments are any other arguments to pass to Terraform
func (terraformer *Terraformer) RunPiped(command string, pipeOutput PipeHandler, arguments ...string) (success bool, err error) {
	var executablePath string
	executablePath, err = terraformer.getExecutablePath()
	if err != nil {
		return
	}

	args := []string{command}
	args = append(args, arguments...)

	commandLine := fmt.Sprintf("%s %s",
		terraformer.ExecutablePath,
		command,
	)
	if len(arguments) > 0 {
		commandLine += " "
		commandLine += strings.Join(arguments, " ")
	}

	var (
		terraformCommand *exec.Cmd
		combinedOutput   io.Reader
		outputCloser     io.Closer
		outputScanner    *bufio.Scanner
	)
	terraformCommand = exec.Command(executablePath, args...)
	combinedOutput, outputCloser, err = pipeCombinedOutput(terraformCommand)
	if err != nil {
		return
	}

	terraformCommand.Dir = terraformer.ConfigDir // Always run Terraform in the configuration directory

	log.Debugf(`Executing "%s" ...`, commandLine)
	err = terraformCommand.Start()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Failed to start: %s", err.Error())
	}

	// Pipe output to the caller.
	outputScanner = bufio.NewScanner(combinedOutput)
	go func() {
		defer outputCloser.Close()

		for outputScanner.Scan() {
			pipeOutput(outputScanner.Text())
		}
	}()

	err = terraformCommand.Wait()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Did not exit cleanly: %s", err.Error())
	}

	if err != nil {
		err = fmt.Errorf("Execute Terraform: Failed (%s)", err.Error())

		return
	}

	success = true

	log.Debugf("Successfully executed %s ...", commandLine)

	return
}

// Get the location of the Terraform executable.
func (terraformer *Terraformer) getExecutablePath() (string, error) {
	err := terraformer.ensureExecutableIsResolved()
	if err != nil {
		return "", err
	}

	return terraformer.ExecutablePath, nil
}

// Determine the location of the Terraform executable.
func (terraformer *Terraformer) resolveExecutablePath() error {
	var err error
	if terraformer.ExecutablePath == "" {
		log.Debugf("Terraform executable location has not been explicitly configured. We will search for it on system PATH.")

		terraformer.ExecutablePath, err = exec.LookPath("terraform")
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Terraform executable location has been explicitly configured.")
		_, err = os.Stat(terraformer.ExecutablePath)
		if err != nil {
			return err
		}
	}

	log.Debugf("The terraformer will use Terraform executable '%s'.", terraformer.ExecutablePath)

	return nil
}

// Ensure that the Terraform executable has been resolved.
func (terraformer *Terraformer) ensureExecutableIsResolved() error {
	if terraformer.ExecutablePath == "" {
		return errors.New("Terraform executable path has not been resolved")
	}

	return nil
}
