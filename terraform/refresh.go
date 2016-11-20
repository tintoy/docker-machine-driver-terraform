package terraform

import (
	"fmt"
	"log"
)

// Refresh invokes Terraform's "refresh" command.
func (terraformer *Terraformer) Refresh() error {
	success, programOutput, err := terraformer.Run("refresh",
		"-input=false", // non-interactive
		"-no-color",
		"-var-file=tfvars.json",
	)
	log.Print(programOutput)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to execute 'terraform refresh'\n:Terraform output:\n%s", programOutput)
	}

	return nil
}
