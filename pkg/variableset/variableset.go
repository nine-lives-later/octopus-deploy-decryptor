package variableset

import (
	"encoding/json"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"os"
)

type VariableSet struct {
	Id        string      `json:"Id"`
	Variables []*Variable `json:"Variables"`
}

type Variable struct {
	Id    string          `json:"Id"`
	Name  string          `json:"Name"`
	Type  string          `json:"Type"`
	Value string          `json:"Value"`
	Scope json.RawMessage `json:"Scope"`
}

func (v *Variable) DecryptedValue(key []byte) (string, error) {
	if v.Type == "Sensitive" {
		return decryptor.DecryptString(key, v.Value)
	}

	return v.Value, nil
}

func ReadVariables(filename string) ([]*Variable, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var set VariableSet
	err = json.Unmarshal(data, &set)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	return set.Variables, nil
}
