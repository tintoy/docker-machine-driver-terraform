package terraform

import (
	"errors"
)

// Destroy invokes Terraform's "destroy" command.
func (terraformer *Terraformer) Destroy(withVariablesFile bool) (success bool, err error) {
	args := []string{
		"-force", "-input=false", // non-interactive
		"-no-color",
	}
	if withVariablesFile {
		args = append(args, "-var-file=tfvars.json")
	}

	success, err = terraformer.RunStreamed("destroy", args...)
	if err != nil {
		return
	}
	if !success {
		err = errors.New("Failed to execute 'terraform destroy'")
	}

	return
}
