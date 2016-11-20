package terraform

import (
	"errors"
)

// Apply invokes Terraform's "apply" command.
func (terraformer *Terraformer) Apply() (success bool, err error) {
	success, err = terraformer.RunStreamed("apply",
		"-input=false", // non-interactive
		"-no-color",
		"-var-file=tfvars.json",
	)
	if err != nil {
		return
	}
	if !success {
		err = errors.New("Failed to execute 'terraform apply'")
	}

	return
}
