package projectExport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type LibraryVariableSet struct {
	SpaceId       string                        `json:"SpaceId"`
	Id            string                        `json:"Id"`
	Name          string                        `json:"Name"`
	VariableSetId string                        `json:"VariableSetId"`
	Templates     []*LibraryVariableSetTemplate `json:"Templates"`
}

func (p *LibraryVariableSet) EntityID() string {
	return p.Id
}

func (p *LibraryVariableSet) EntityName() string {
	return p.Name
}

func (p *LibraryVariableSet) EntitySpaceID() string {
	return p.SpaceId
}

func (p *LibraryVariableSet) AddToEntityMap(m EntityMap) {
	m[p.EntityID()] = p

	for _, t := range p.Templates {
		t.AddToEntityMap(m)
	}
}

type LibraryVariableSetTemplate struct {
	Id           string        `json:"Id"`
	Name         string        `json:"Name"`
	DefaultValue *DefaultValue `json:"DefaultValue"`

	// injected by `ReadProject()`
	SpaceId string `json:"-"`
}

func (t *LibraryVariableSetTemplate) EntityID() string {
	return t.Id
}

func (t *LibraryVariableSetTemplate) EntityName() string {
	return t.Name
}

func (t *LibraryVariableSetTemplate) EntitySpaceID() string {
	return t.SpaceId
}

func (t *LibraryVariableSetTemplate) AddToEntityMap(m EntityMap) {
	m[t.EntityID()] = t
}

func (t *LibraryVariableSetTemplate) DecryptedValue(key []byte) (string, error) {
	return t.DefaultValue.DecryptedValue(key)
}

var (
	_ Entity = &LibraryVariableSet{}         // ensure type compatibility
	_ Entity = &LibraryVariableSetTemplate{} // ensure type compatibility
)

func IsLibraryVariableSetFilename(f string) bool {
	return strings.HasPrefix(f, "LibraryVariableSets-") && strings.HasSuffix(f, ".json")
}

func ReadLibraryVariableSet(filename string) (*LibraryVariableSet, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret LibraryVariableSet
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	// inject space into templates
	for _, v := range ret.Templates {
		v.SpaceId = ret.SpaceId
	}

	return &ret, nil
}
