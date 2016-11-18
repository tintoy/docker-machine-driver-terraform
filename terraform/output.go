package terraform

import (
	"encoding/json"
	"fmt"
	"log"
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

// Invoke Terraform's "output" command and parse the results.
//
// Returns a map of outputs, keyed by name.
func (terraformer *Terraformer) runTerraformOutput() (success bool, outputs Outputs, err error) {
	var programOutput string
	success, programOutput, err = terraformer.Run("output",
		"-json",
	)
	log.Print(programOutput)

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
