package main

/*
 * Driver implementation (Terraform configuration)
 * -----------------------------------------------
 */

import (
	"fmt"
	"os"
	"path"

	"github.com/docker/machine/libmachine/log"
)

func (driver *Driver) getVariablesFileName() (string, error) {
	localConfigDir, err := driver.getConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(localConfigDir, "tfvars.json"), nil
}

// Read Terraform variables from tfvars.json (if it exists)
func (driver *Driver) readVariables() error {
	variablesFileName, err := driver.getVariablesFileName()
	if err != nil {
		return err
	}

	_, err = os.Stat(variablesFileName)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	log.Debugf("Reading Terraform variables from '%s'...",
		variablesFileName,
	)

	driver.ConfigVariables.Clear()
	err = driver.ConfigVariables.Read(variablesFileName)
	if err != nil {
		return err
	}

	return nil
}

// Read Terraform variables from the additional variables file passed in on the command-line.
func (driver *Driver) readAdditionalVariables() error {
	if driver.AdditionalVariablesFile == "" {
		return nil // Nothing to do
	}

	log.Debugf("Reading additional Terraform variables from '%s'...",
		driver.AdditionalVariablesFile,
	)

	// This operation is additive (preserves exising variables).
	err := driver.ConfigVariables.Read(driver.AdditionalVariablesFile)
	if err != nil {
		return fmt.Errorf("Unable to read additional variables from '%s': %s",
			driver.AdditionalVariablesFile,
			err.Error(),
		)
	}

	return nil
}

// Write Terraform variables to tfvars.json
func (driver *Driver) writeVariables() error {
	variablesFileName, err := driver.getVariablesFileName()
	if err != nil {
		return err
	}

	log.Debugf("Writing %d Terraform variables to '%s'...",
		len(driver.ConfigVariables),
		variablesFileName,
	)

	err = driver.ConfigVariables.Write(variablesFileName)
	if err != nil {
		return err
	}

	return nil
}
