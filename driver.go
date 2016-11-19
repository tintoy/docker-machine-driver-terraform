package main

/*
 * Driver implementation
 * ---------------------
 */

import (
	"errors"
	"os"

	stdlog "log"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	"github.com/tintoy/docker-machine-driver-terraform/terraform"
)

// Driver is the Docker Machine driver for Terraform.
type Driver struct {
	*drivers.BaseDriver

	// The source path (or URL) of the Terraform configuration.
	ConfigSource string

	// The path of the directory containing the imported Terraform configuration.
	ConfigDir string

	// Additional variables for the Terraform configuration
	ConfigVariables terraform.ConfigVariables

	// An optional file containing the JSON that represents additional variables for the Terraform configuration
	AdditionalVariablesFile string

	// Refresh the configuration after applying it
	RefreshAfterApply bool

	// The full path to the Terraform executable.
	TerraformExecutablePath string

	// The path to the SSH private key file to use for authentication.
	//
	// If not specified, a new key-pair will be generated.
	SSHKey string

	// The terraform executor.
	terraformer *terraform.Terraformer
}

// GetCreateFlags registers the "machine create" flags recognized by this driver, including
// their help text and defaults.
func (driver *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:  "terraform-config-source",
			Usage: "The path (or URL) of the Terraform configuration",
			Value: "",
		},
		mcnflag.StringFlag{
			Name:  "terraform-additional-variables",
			Usage: "An optional file containing the JSON that represents additional variables for the Terraform configuration",
			Value: "",
		},
		mcnflag.BoolFlag{
			Name:  "terraform-refresh",
			Usage: "Refresh the configuration after applying it",
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
	log.Debugf("docker-machine-driver-terraform %s", DriverVersion)

	// Enable ALL logging if MACHINE_DEBUG is set
	if os.Getenv("MACHINE_DEBUG") != "" {
		stdlog.SetOutput(os.Stderr)
	}

	driver.ConfigSource = flags.String("terraform-config-source")
	driver.ConfigVariables = make(map[string]interface{})
	driver.AdditionalVariablesFile = flags.String("terraform-additional-variables")
	driver.RefreshAfterApply = flags.Bool("terraform-refresh")

	driver.SSHPort = flags.Int("terraform-ssh-port")
	driver.SSHUser = flags.String("terraform-ssh-user")
	driver.SSHKey = flags.String("terraform-ssh-key")

	// Validation
	if driver.ConfigSource == "" {
		return errors.New("Required argument: --terraform-config-source")
	}

	return nil
}

// PreCreateCheck validates the configuration before making any changes.
func (driver *Driver) PreCreateCheck() error {
	if driver.ConfigSource == "" {
		return errors.New("The source for Terraform configuration has not been specified")
	}

	log.Infof("Will create machine '%s' using Terraform configuration from '%s'.",
		driver.MachineName,
		driver.ConfigSource,
	)

	log.Infof("Resolving Terraform configuration...")
	err := driver.resolveConfigDir()
	if err != nil {
		return err
	}
	err = driver.importConfig()
	if err != nil {
		return err
	}

	log.Infof("Customising terraform configuration...")
	driver.ConfigVariables["dm_machine_name"] = driver.MachineName
	driver.ConfigVariables["dm_ssh_user"] = driver.SSHUser
	driver.ConfigVariables["dm_ssh_port"] = driver.SSHPort
	err = driver.readAdditionalVariables()
	if err != nil {
		return err
	}
	err = driver.writeVariables()
	if err != nil {
		return err
	}

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
