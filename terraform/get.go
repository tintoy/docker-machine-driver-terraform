package terraform

import (
	"fmt"
	"log"
)

// Get invokes Terraform's "get" command and parse the results.
//
// Returns a map of outputs, keyed by name.
func (terraformer *Terraformer) Get() error {
	success, programOutput, err := terraformer.Run("get",
		"-no-color",
	)
	log.Print(programOutput)

	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("Failed to execute 'terraform get'")
	}

	return nil
}
