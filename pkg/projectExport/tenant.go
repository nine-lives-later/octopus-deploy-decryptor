package projectExport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Tenant struct {
	Id      string `json:"Id"`
	Name    string `json:"Name"`
	SpaceId string `json:"SpaceId"`
}

var (
	_ Entity = &Tenant{} // ensure type compatibility
)

func (t *Tenant) EntityID() string {
	return t.Id
}

func (t *Tenant) EntityName() string {
	return t.Name
}

func (t *Tenant) EntitySpaceID() string {
	return t.SpaceId
}

func (t *Tenant) AddToEntityMap(m EntityMap) {
	m[t.EntityID()] = t
}

func IsTenantFilename(f string) bool {
	return strings.HasPrefix(f, "Tenants-") && strings.HasSuffix(f, ".json")
}

func ReadTenant(filename string) (*Tenant, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret Tenant
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	return &ret, nil
}
