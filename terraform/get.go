package terraform

import (
	"fmt"
	"log"
)

// Get invokes Terraform's "get" command.
func (terraformer *Terraformer) Get() error {
	success, programOutput, err := terraformer.Run("get",
		"-no-color",
	)
	log.Print(programOutput)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to execute 'terraform get'\n:Terraform output:\n%s", programOutput)
	}

	return nil
}
