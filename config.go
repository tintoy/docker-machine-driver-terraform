package main

/*
 * Driver implementation (Terraform configuration)
 * -----------------------------------------------
 */

import (
	"errors"

	"github.com/hashicorp/go-getter"
)

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
