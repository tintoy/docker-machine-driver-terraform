package terraform

import (
	"fmt"
	"log"
)

// Validate invokes Terraform's "validate" command and parse the results.
//
// Returns a map of outputs, keyed by name.
func (terraformer *Terraformer) Validate() error {
	success, programOutput, err := terraformer.Run("validate")
	log.Print(programOutput)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to execute 'terraform validate'\n:Terraform output:\n%s", programOutput)
	}

	return nil
}
