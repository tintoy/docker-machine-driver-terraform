package main

/*
 * Driver implementation (Terraform configuration)
 * -----------------------------------------------
 */

import (
	"fmt"
	"os"
	"path"
	"strings"

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

// Read all additional variables (from the command line, or from a file)
func (driver *Driver) readAdditionalVariables() error {
	err := driver.readAdditionalVariablesInline()
	if err != nil {
		return err
	}

	return driver.readAdditionalVariablesFile()
}

// Read additional variables passed in on the command-line (--terraform-variable a=b --terraform-variable c=d)
func (driver *Driver) readAdditionalVariablesInline() error {
	for _, additionalVariable := range driver.AdditionalVariablesInline {
		variableNameAndValue := strings.SplitN(additionalVariable, "=", 2)
		if len(variableNameAndValue) != 2 {
			return fmt.Errorf("Invalid format for additional variable '%s", additionalVariable)
		}

		// Don't overwrite existing entries.
		_, variableExists := driver.ConfigVariables[variableNameAndValue[0]]
		if variableExists {
			continue
		}

		driver.ConfigVariables[variableNameAndValue[0]] = variableNameAndValue[1]
	}

	return nil
}

// Read Terraform variables from the file passed in on the command-line.
func (driver *Driver) readAdditionalVariablesFile() error {
	if driver.AdditionalVariablesFile == "" {
		return nil // Nothing to do
	}

	variablesFileName := driver.AdditionalVariablesFile
	if !path.IsAbs(variablesFileName) {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return err
		}
		variablesFileName = path.Join(workingDirectory, variablesFileName)
	}

	log.Debugf("Reading additional Terraform variables from '%s'...", variablesFileName)

	// This operation is additive (preserves exising variables).
	err := driver.ConfigVariables.Read(variablesFileName)
	if err != nil {
		return fmt.Errorf("Unable to read additional variables from '%s': %s",
			variablesFileName,
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
