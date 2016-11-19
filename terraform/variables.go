package terraform

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
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
	normalizedVariables := variables.normalize()

	variablesJSON, err := json.MarshalIndent(normalizedVariables, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, variablesJSON, 0644 /* u=rw,g=r,o=r */)
	if err != nil {
		return err
	}

	return nil
}

// Make a copy of the configuration variables, but with normalised values.
//
// For example, numbers are converted to strings.
func (variables ConfigVariables) normalize() ConfigVariables {
	normalized := make(ConfigVariables)
	for variableName, variableValue := range variables {
		normalizedValue := variableValue

		integerValue, ok := variableValue.(int)
		if ok {
			normalizedValue = strconv.Itoa(integerValue)
		}

		normalized[variableName] = normalizedValue
	}

	return normalized
}
