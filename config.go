package main

/*
 * Driver implementation (Terraform configuration)
 * -----------------------------------------------
 */

import (
	"errors"

	"github.com/docker/machine/libmachine/log"
	"github.com/tintoy/docker-machine-driver-terraform/fetch"
)

func (driver *Driver) importConfig() error {
	terraformer, err := driver.getTerraformer()
	if err != nil {
		return err
	}

	// Directory won't exist until fetch (go-getter) creates it.
	localConfigDir, err := driver.getConfigDir()
	if err != nil {
		return err
	}

	log.Infof("Importing Terraform configuration from '%s' to '%s'...",
		driver.ConfigSource,
		localConfigDir,
	)
	driver.ConfigSource, err = fetch.ParseSource(driver.ConfigSource)
	if err != nil {
		return err
	}

	log.Debugf("Fetching Terraform configuration from '%s...'", driver.ConfigSource)
	err = fetch.Content(driver.ConfigSource, localConfigDir)
	if err != nil {
		return err
	}

	log.Debugf("Fetching Terraform modules (if any)...")
	err = terraformer.Get()
	if err != nil {
		return err
	}

	log.Infof("Validating Terraform configuration...")
	err = terraformer.Validate()
	if err != nil {
		return err
	}

	log.Debugf("Import complete.")

	return nil
}

// Get the local directory where the resolved Terraform configuration is cached.
func (driver *Driver) getConfigDir() (string, error) {
	err := driver.ensureConfigDirIsResolved()
	if err != nil {
		return "", err
	}

	return driver.ConfigDir, nil
}

// Ensure that the local Terraform configuration directory exists.
func (driver *Driver) resolveConfigDir() error {
	if driver.ConfigDir != "" {
		return errors.New("Local Terraform configuration directory has already been resolved")
	}

	// We don't auto-create this directory because go-getter wants to do it for us.
	driver.ConfigDir = driver.ResolveStorePath("terraform-config")

	return nil
}

// Ensure that the Terraform configuration has been resolved.
func (driver *Driver) ensureConfigDirIsResolved() error {
	if driver.ConfigDir == "" {
		return errors.New("Terraform configuration has not been resolved")
	}

	return nil
}
