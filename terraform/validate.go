package terraform

import (
	"errors"

	"github.com/docker/machine/libmachine/log"
)

// Validate invokes Terraform's "validate" command and parse the results.
//
// Returns a map of outputs, keyed by name.
func (terraformer *Terraformer) Validate() error {
	success, programOutput, err := terraformer.Run("validate")
	log.Info(programOutput)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("Failed to execute 'terraform validate'")
	}

	return nil
}
