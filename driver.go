package main

/*
 * Driver implementation
 * ---------------------
 */

import (
	"errors"
	"fmt"
	"net"
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

	// Optional "name=value" items that represent additional variables for the Terraform configuration
	AdditionalVariablesInline []string

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
			Name:  "terraform-config",
			Usage: "The path (or URL) of the Terraform configuration",
			Value: "",
		},
		mcnflag.StringSliceFlag{
			Name:  "terraform-variable",
			Usage: "Additional variable(s) for the Terraform configuration (in the form name=value)",
			Value: []string{},
		},
		mcnflag.StringFlag{
			Name:  "terraform-variables-from",
			Usage: "The name of a file containing the JSON that represents additional variables for the Terraform configuration",
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

	driver.ConfigSource = flags.String("terraform-config")
	driver.ConfigVariables = make(map[string]interface{})

	driver.AdditionalVariablesInline = flags.StringSlice("terraform-variable")
	driver.AdditionalVariablesFile = flags.String("terraform-variables")

	driver.RefreshAfterApply = flags.Bool("terraform-refresh")

	driver.SSHPort = flags.Int("terraform-ssh-port")
	driver.SSHUser = flags.String("terraform-ssh-user")
	driver.SSHKey = flags.String("terraform-ssh-key")

	// Validation
	if driver.ConfigSource == "" {
		return errors.New("Required argument: --terraform-config")
	}

	return nil
}

// PreCreateCheck validates the configuration before making any changes.
func (driver *Driver) PreCreateCheck() error {
	if driver.ConfigSource == "" {
		return errors.New("The source for Terraform configuration has not been specified")
	}

	log.Infof("Auto-detecting client's public (external) IP address...")
	clientIP, err := getClientPublicIPv4Address()
	if err != nil {
		return err
	}

	log.Infof("Will create machine '%s' using Terraform configuration from '%s'.",
		driver.MachineName,
		driver.ConfigSource,
	)

	log.Infof("Resolving Terraform configuration...")
	err = driver.resolveConfigDir()
	if err != nil {
		return err
	}
	err = driver.importConfig()
	if err != nil {
		return err
	}

	if driver.SSHKey != "" {
		log.Infof("Importing SSH key '%s'...", driver.SSHKey)
		err = driver.importSSHKey()
		if err != nil {
			return err
		}
	} else {
		log.Infof("Generating new SSH key...")
		err = driver.generateSSHKey()
		if err != nil {
			return err
		}
	}

	log.Infof("Customising terraform configuration...")
	driver.ConfigVariables["dm_client_ip"] = clientIP
	driver.ConfigVariables["dm_machine_name"] = driver.MachineName
	driver.ConfigVariables["dm_ssh_private_key_file"] = driver.SSHKeyPath
	driver.ConfigVariables["dm_ssh_public_key_file"] = driver.SSHKeyPath + ".pub"
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
	log.Infof("Applying Terraform configuration...")

	terraformer, err := driver.getTerraformer()
	if err != nil {
		return err
	}

	success, err := terraformer.Apply()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("Failed to apply Terraform configuration")
	}

	if driver.RefreshAfterApply {
		log.Infof("Refreshing Terraform configuration state...")
		err = terraformer.Refresh()
		if err != nil {
			return err
		}
	}

	outputs, err := terraformer.Output()
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to obtain Terraform outputs")
	}

	output, ok := outputs["dm_machine_ip"]
	if !ok {
		return fmt.Errorf("Configuration does not declare required output 'dm_machine_ip'")
	}
	driver.IPAddress = output.Value.(string)

	output, ok = outputs["dm_ssh_user"]
	if ok {
		driver.SSHUser = output.Value.(string)
	}

	log.Infof("Deployed host has IP '%s'.", driver.IPAddress)
	log.Infof("Deployed host has SSH user '%s'.", driver.SSHUser)

	return nil
}

// GetState retrieves the status of the target Docker Machine instance in CloudControl.
func (driver *Driver) GetState() (state.State, error) {
	return state.Running, nil
}

// GetURL returns docker daemon URL on the target machine
func (driver *Driver) GetURL() (string, error) {
	if driver.IPAddress == "" {
		return "", nil
	}

	url := fmt.Sprintf("tcp://%s", net.JoinHostPort(driver.IPAddress, "2376"))

	return url, nil
}

// Remove deletes the target machine.
func (driver *Driver) Remove() error {
	log.Infof("Destroying terraform configuration...")

	terraformer, err := driver.getTerraformer()
	if err != nil {
		return err
	}

	success, err := terraformer.Destroy()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("Failed to destroy Terraform configuration")
	}

	return nil
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
