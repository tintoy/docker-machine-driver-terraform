package main

/*
 * Driver implementation (Terraform integration)
 * ---------------------------------------------
 */

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/machine/libmachine/log"
)

// Represents the result of Terraform's "output" command.
type terraformOutput struct {
	DataType  string      `json:"type"`
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
}

// A map of Terraform outputs, keyed by name.
type terraformOutputs map[string]terraformOutput

// Invoke Terraform's "output" command and parses the results.
//
// Returns a map of outputs, keyed by name.
func (driver *Driver) runTerraformOutput() (success bool, outputs terraformOutputs, err error) {
	var programOutput string
	success, programOutput, err = driver.runTerraform("output",
		"-json",
	)
	log.Debug(programOutput)

	outputs = make(terraformOutputs)
	err = json.Unmarshal(
		[]byte(programOutput),
		&outputs,
	)
	if err != nil {
		err = fmt.Errorf("Failed to parse JSON from Terraform output: %s ", err.Error())

		return
	}

	return
}

// Invoke Terraform.
//
// command is the name of the Terraform command to execute (e.g. plan, apply, output,  destroy, etc)
// arguments are any other arguments to pass to Terraform
func (driver *Driver) runTerraform(command string, arguments ...string) (success bool, output string, err error) {
	err = driver.ensureTerraformExecutableIsResolved()
	if err != nil {
		return
	}

	args := []string{command}
	args = append(args, arguments...)

	commandLine := fmt.Sprintf(`"%s" %s`,
		driver.TerraformExecutablePath,
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
	terraformCommand = exec.Command(command, args...)
	terraformCommand.Stdout = &programOutput
	terraformCommand.Stderr = &programOutput
	terraformCommand.Dir = driver.ConfigDir // Always run Terraform in the cached configuration directory

	log.Debugf("Executing %s ...", commandLine)
	err = terraformCommand.Run()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Failed: %s", err.Error())

		output = string(
			programOutput.Bytes(),
		)

		return
	}

	err = terraformCommand.Wait()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Process did not exit cleanly: %s", err.Error())

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

// Get the location of the Terraform executable.
func (driver *Driver) getTerraformExecutablePath() (executablePath string, err error) {
	err = driver.ensureTerraformExecutableIsResolved()
	if err != nil {
		return
	}

	executablePath = driver.TerraformExecutablePath

	return
}

// Determine the location of the Terraform executable.
func (driver *Driver) resolveTerraformExecutablePath() error {
	var err error
	if driver.TerraformExecutablePath == "" {
		log.Debugf("Terraform executable location has not been explicitly configured. We will search for it on system PATH.")
		driver.TerraformExecutablePath, err = exec.LookPath("terraform")
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Terraform executable location has been explicitly configured.")
		_, err = os.Stat(driver.TerraformExecutablePath)
		if err != nil {
			return err
		}
	}

	log.Debugf("The driver will use Terraform executable '%s'.", driver.TerraformExecutablePath)

	return nil
}

// Ensure that the Terraform executable has been resolved.
func (driver *Driver) ensureTerraformExecutableIsResolved() error {
	if driver.TerraformExecutablePath == "" {
		return errors.New("Terraform executable path has not been resolved")
	}

	return nil
}

// Get the local directory where the resolved Terraform configuration is cached.
func (driver *Driver) getConfigDir() (configDir string, err error) {
	err = driver.ensureConfigDirIsResolved()
	if err != nil {
		return
	}

	configDir = driver.ConfigDir

	return
}

// Locate and cache the Terraform configuration.
func (driver *Driver) resolveConfigDir() error {
	if driver.ConfigDir != "" {
		return errors.New("Terraform configuration has already been resolved")
	}

	// TODO: Support other configuration sources besides local directory
	driver.ConfigDir = driver.ConfigSource

	return nil
}

// Ensure that the Terraform configuration has been resolved.
func (driver *Driver) ensureConfigDirIsResolved() error {
	if driver.ConfigDir == "" {
		return errors.New("Terraform configuration has not been resolved")
	}

	return nil
}
