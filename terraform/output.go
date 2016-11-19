package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/docker/machine/libmachine/log"
)

// Output represents an output from Terraform's "output" command.
type Output struct {
	Name      string      `json:""`
	DataType  string      `json:"type"`
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive"`
}

// Outputs is a map of Terraform outputs, keyed by name.
type Outputs map[string]Output

// Output invokes Terraform's "output" command and parse the results.
//
// Returns a map of outputs, keyed by name.
func (terraformer *Terraformer) Output() (success bool, outputs Outputs, err error) {
	var programOutput string
	success, programOutput, err = terraformer.Run("output",
		"-json",
	)
	log.Debugf(programOutput)
	if err != nil {
		return
	}
	if !success {
		err = fmt.Errorf("Failed to execute 'terraform output'\n:Terraform output:\n%s", programOutput)

		return
	}

	outputs = make(Outputs)
	err = json.Unmarshal(
		[]byte(programOutput),
		&outputs,
	)
	if err != nil {
		err = fmt.Errorf("Failed to parse JSON from Terraform output: %s ", err.Error())

		return
	}

	return
}
