package terraform

import (
	"errors"
	"fmt"

	"github.com/docker/machine/libmachine/log"
)

// Apply invokes Terraform's "apply" command.
func (terraformer *Terraformer) Apply(variablesFilePath string) (success bool, programOutput string, err error) {
	args := []string{
		"-input=false", // non-interactive
		"-no-color",
	}
	if variablesFilePath != "" {
		args = append(args,
			fmt.Sprintf("-var-file=%s", variablesFilePath),
		)
	}

	pipeHandler := func(outputLine string) {
		log.Debug(outputLine)
	}

	success, err = terraformer.RunPiped("apply", pipeHandler, args...)
	if err != nil {
		return
	}
	if !success {
		err = errors.New("Failed to execute 'terraform apply'")
	}

	return
}
