package terraform

import (
	"errors"
)

// Get invokes Terraform's "get" command.
func (terraformer *Terraformer) Get() error {
	success, err := terraformer.RunStreamed("get", "-no-color")
	if err != nil {
		return err
	}
	if !success {
		return errors.New("Failed to execute 'terraform get'")
	}

	return nil
}
