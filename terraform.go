package main

/*
 * Driver implementation (Terraform integration)
 * ---------------------------------------------
 */

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/machine/libmachine/log"
)

// runTerraform executes a Terraform command.
// command is the name of the command to execute (e.g. plan, apply, output,  destroy, etc)
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

	log.Debugf("Executing %s ...", commandLine)
	err = terraformCommand.Run()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Failed: %s", err.Error())

		return
	}

	err = terraformCommand.Wait()
	if err != nil {
		err = fmt.Errorf("Execute Terraform: Process did not exit cleanly: %s", err.Error())

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

func (driver *Driver) ensureTerraformExecutableIsResolved() error {
	if driver.TerraformExecutablePath == "" {
		return errors.New("Terraform executable path has not been resolved")
	}

	return nil
}
