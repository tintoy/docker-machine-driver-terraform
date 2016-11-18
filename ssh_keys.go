package main

/*
 * SSH key generation / import
 * ---------------------------
 */

import (
	"errors"
	"fmt"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"io/ioutil"
	"os"
	"path"
)

// Generate an SSH key pair, and save it into the machine store folder.
func (driver *Driver) generateSSHKey() error {
	if driver.SSHKeyPath != "" {
		return errors.New("SSH key path already configured")
	}

	driver.SSHKeyPath = driver.ResolveStorePath("id_rsa")
	err := ssh.GenerateSSHKey(driver.SSHKeyPath)
	if err != nil {
		log.Errorf("Failed to generate SSH key pair: %s", err.Error())

		return err
	}

	return nil
}

// Import the configured SSH key files into the machine store folder.
func (driver *Driver) importSSHKey() error {
	if driver.SSHKey == "" {
		return errors.New("SSH key path not configured")
	}

	driver.SSHKeyPath = driver.ResolveStorePath(
		path.Base(driver.SSHKey),
	)
	err := copySSHKey(driver.SSHKey, driver.SSHKeyPath)
	if err != nil {
		log.Infof("Couldn't copy SSH private key: %s", err.Error())

		return err
	}

	err = copySSHKey(driver.SSHKey+".pub", driver.SSHKeyPath+".pub")
	if err != nil {
		log.Infof("Couldn't copy SSH public key: %s", err.Error())

		return err
	}

	return nil
}

// Get the public portion of the configured SSH key.
func (driver *Driver) getSSHPublicKey() (string, error) {
	publicKeyFile, err := os.Open(driver.SSHKeyPath + ".pub")
	if err != nil {
		return "", err
	}
	defer publicKeyFile.Close()

	publicKeyData, err := ioutil.ReadAll(publicKeyFile)
	if err != nil {
		return "", err
	}

	return string(publicKeyData), nil
}

// Copy an SSH key file.
func copySSHKey(sourceFile string, destinationFile string) error {
	err := mcnutils.CopyFile(sourceFile, destinationFile)
	if err != nil {
		return fmt.Errorf("unable to copy ssh key: %s", err.Error())
	}

	err = os.Chmod(destinationFile, 0600)
	if err != nil {
		return fmt.Errorf("unable to set permissions on the ssh key: %s", err.Error())
	}

	return nil
}
