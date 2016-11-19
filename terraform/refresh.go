package terraform

import (
	"fmt"
	"log"
)

// Refresh invokes Terraform's "refresh" command.
func (terraformer *Terraformer) Refresh(variablesFilePath string) error {
	args := []string{
		"-input=false", // non-interactive
		"-no-color",
	}
	if variablesFilePath != "" {
		args = append(args,
			fmt.Sprintf("-var-file=%s", variablesFilePath),
		)
	}

	success, programOutput, err := terraformer.Run("refresh", args...)
	log.Print(programOutput)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("Failed to execute 'terraform refresh'\n:Terraform output:\n%s", programOutput)
	}

	return nil
}
