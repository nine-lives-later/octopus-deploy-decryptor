package projectExport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Environment struct {
	Type    string `json:"Type"`
	Id      string `json:"Id"`
	Name    string `json:"Name"`
	SpaceId string `json:"SpaceId"`
}

var (
	_ Entity = &Environment{} // ensure type compatibility
)

func (e *Environment) EntityID() string {
	return e.Id
}

func (e *Environment) EntityName() string {
	return e.Name
}

func (e *Environment) EntitySpaceID() string {
	return e.SpaceId
}

func (e *Environment) AddToEntityMap(m EntityMap) {
	m[e.EntityID()] = e
}

func IsEnvironmentFilename(f string) bool {
	return strings.HasPrefix(f, "Environments-") && strings.HasSuffix(f, ".json")
}

func ReadEnvironment(filename string) (*Environment, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret Environment
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	return &ret, nil
}
