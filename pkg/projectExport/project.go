package projectExport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Project struct {
	SpaceId                       string             `json:"SpaceId"`
	Id                            string             `json:"Id"`
	Name                          string             `json:"Name"`
	Slug                          string             `json:"Slug"`
	VariableSetId                 string             `json:"VariableSetId"`
	IncludedLibraryVariableSetIds []string           `json:"IncludedLibraryVariableSetIds"`
	Templates                     []*ProjectTemplate `json:"Templates"`
}

func (p *Project) EntityID() string {
	return p.Id
}

func (p *Project) EntityName() string {
	return p.Name
}

func (p *Project) EntitySpaceID() string {
	return p.SpaceId
}

func (p *Project) AddToEntityMap(m EntityMap) {
	m[p.EntityID()] = p

	for _, t := range p.Templates {
		t.AddToEntityMap(m)
	}
}

type ProjectTemplate struct {
	Id           string        `json:"Id"`
	Name         string        `json:"Name"`
	DefaultValue *DefaultValue `json:"DefaultValue"`

	// injected by `ReadProject()`
	SpaceId string `json:"-"`
}

func (t *ProjectTemplate) EntityID() string {
	return t.Id
}

func (t *ProjectTemplate) EntityName() string {
	return t.Name
}

func (t *ProjectTemplate) EntitySpaceID() string {
	return t.SpaceId
}

func (t *ProjectTemplate) AddToEntityMap(m EntityMap) {
	m[t.EntityID()] = t
}

func (t *ProjectTemplate) DecryptedValue(key []byte) (string, error) {
	return t.DefaultValue.DecryptedValue(key)
}

var (
	_ Entity = &Project{}         // ensure type compatibility
	_ Entity = &ProjectTemplate{} // ensure type compatibility
)

func IsProjectFilename(f string) bool {
	return strings.HasPrefix(f, "Projects-") && strings.HasSuffix(f, ".json")
}

func ReadProject(filename string) (*Project, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret Project
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
