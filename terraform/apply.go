package terraform

import (
	"errors"
)

// Apply invokes Terraform's "apply" command.
func (terraformer *Terraformer) Apply(withVariablesFile bool) (success bool, err error) {
	args := []string{
		"-input=false", // non-interactive
		"-no-color",
	}
	if withVariablesFile {
		args = append(args, "-var-file=tfvars.json")
	}

	success, err = terraformer.RunStreamed("apply", args...)
	if err != nil {
		return
	}
	if !success {
		err = errors.New("Failed to execute 'terraform apply'")
	}

	return
}
