package main

/*
 * Driver implementation (Terraform integration)
 * ---------------------------------------------
 */

import (
	"errors"
	"os"
	"os/exec"

	"github.com/docker/machine/libmachine/log"
)

// Get the location of the Terraform executable.
func (driver *Driver) getTerraformExecutablePath() (executablePath string, err error) {
	if driver.TerraformExecutablePath == "" {
		err = errors.New("Terraform executable path has not been resolved")

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
