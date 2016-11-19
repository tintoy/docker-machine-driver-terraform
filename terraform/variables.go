package terraform

import (
	"encoding/json"
	"io/ioutil"
)

// ConfigVariables represents values for Terraform configuration variables, keyed by name.
type ConfigVariables map[string]interface{}

// Clear the contents of the variable map.
func (variables ConfigVariables) Clear() {
	for variableName := range variables {
		delete(variables, variableName)
	}
}

// Read variables from the specified file (normally tfvars.json).
//
// This operation is additive but does not overwrite existing values.
// Call Clear before Read if this isn't what you want.
func (variables ConfigVariables) Read(fileName string) error {
	variablesJSON, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	var variablesFromJSON ConfigVariables
	err = json.Unmarshal(variablesJSON, &variablesFromJSON)
	if err != nil {
		return err
	}

	for variableName, variableValue := range variablesFromJSON {
		_, ok := variables[variableName]
		if ok {
			continue // Don't overwrite existing values
		}
		variables[variableName] = variableValue
	}

	return nil
}

// Save variables from the specified file (normally tfvars.json).
func (variables ConfigVariables) Write(fileName string) error {
	variablesJSON, err := json.MarshalIndent(variables, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, variablesJSON, 0644 /* u=rw,g=r,o=r */)
	if err != nil {
		return err
	}

	return nil
}
