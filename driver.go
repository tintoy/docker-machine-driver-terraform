package main

/*
 * Driver implementation
 * ---------------------
 */

import (
	"errors"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

// Driver is the Docker Machine driver for Terraform.
type Driver struct {
	*drivers.BaseDriver

	// The path (or URL) of the Terraform configuration.
	ConfigLocation string
}

// GetCreateFlags registers the "machine create" flags recognized by this driver, including
// their help text and defaults.
func (driver *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "TERRAFORM_CONFIG",
			Name:   "terraform-config",
			Usage:  "The path (or URL) of the Terraform configuration",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "TERRAFORM_SSH_USER",
			Name:   "terraform-ssh-user",
			Usage:  "The SSH username to use. Default: root",
			Value:  "root",
		},
		mcnflag.StringFlag{
			EnvVar: "TERRAFORM_SSH_KEY",
			Name:   "terraform-ssh-key",
			Usage:  "The SSH key file to use",
			Value:  "",
		},
		mcnflag.IntFlag{
			EnvVar: "TERRAFORM_SSH_PORT",
			Name:   "terraform-ssh-port",
			Usage:  "The SSH port. Default: 22",
			Value:  22,
		},
	}
}

// DriverName returns the name of the driver
func (driver *Driver) DriverName() string {
	return "terraform"
}

// SetConfigFromFlags assigns and verifies the command-line arguments presented to the driver.
func (driver *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	driver.ConfigLocation = flags.String("terraform-config")

	driver.SSHPort = flags.Int("terraform-ssh-port")
	driver.SSHUser = flags.String("terraform-ssh-user")
	driver.SSHKeyPath = flags.String("terraform-ssh-key")

	log.Debugf("docker-machine-driver-terraform %s", DriverVersion)

	return nil
}

// PreCreateCheck validates the configuration before making any changes.
func (driver *Driver) PreCreateCheck() error {
	if driver.ConfigLocation == "" {
		return errors.New("")
	}

	log.Infof("Will create machine '%s' using Terraform configuration from '%s'.",
		driver.MachineName,
		driver.ConfigLocation,
	)

	log.Infof("Resolving Terraform configuration...")

	// TODO: Fetch and / or validate configuration as required.

	return nil
}

// Create a new Docker Machine instance on CloudControl.
func (driver *Driver) Create() error {
	return errors.New("Create is not yet implemented.")
}

// GetState retrieves the status of the target Docker Machine instance in CloudControl.
func (driver *Driver) GetState() (state.State, error) {
	return state.None, errors.New("GetState is not yet implemented.")
}

// GetURL returns docker daemon URL on the target machine
func (driver *Driver) GetURL() (string, error) {
	return "", errors.New("GetURL is not yet implemented.")
}

// Remove deletes the target machine.
func (driver *Driver) Remove() error {
	return errors.New("Remove is not yet implemented.")
}

// Start the target machine.
func (driver *Driver) Start() error {
	return errors.New("The Terraform driver does not support Start.")
}

// Stop the target machine (gracefully).
func (driver *Driver) Stop() error {
	return errors.New("The Terraform driver does not support Stop.")
}

// Restart the target machine.
func (driver *Driver) Restart() error {
	// TODO: Check machine has been created.

	_, err := drivers.RunSSHCommandFromDriver(driver, "sudo shutdown -r now")

	return err
}

// Kill the target machine (hard shutdown).
func (driver *Driver) Kill() error {
	return errors.New("The Terraform driver does not support Kill.")
}

// GetSSHHostname returns the hostname for SSH
func (driver *Driver) GetSSHHostname() (string, error) {
	// TODO: Check machine has been created.

	return driver.IPAddress, nil
}

// GetSSHKeyPath returns the ssh key path
func (driver *Driver) GetSSHKeyPath() string {
	return driver.SSHKeyPath
}
