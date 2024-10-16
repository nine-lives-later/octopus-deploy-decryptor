package projectExport

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TenantVariable struct {
	Id                 string       `json:"Id"`
	TenantId           string       `json:"TenantId"`
	OwnerId            string       `json:"OwnerId"`
	OwnerType          string       `json:"OwnerType"`
	EnvironmentId      string       `json:"EnvironmentId"`
	VariableTemplateId string       `json:"VariableTemplateId"`
	SpaceId            string       `json:"SpaceId"`
	Value              DefaultValue `json:"Value"`

	// allow for resolving for variable name
	entityMap EntityMap
}

var (
	_ Entity = &TenantVariable{} // ensure type compatibility
)

func (v *TenantVariable) EntityID() string {
	return v.Id
}

func (v *TenantVariable) EntityName() string {
	// resolve the variable template
	if v.entityMap != nil {
		if n := v.entityMap[v.VariableTemplateId]; n != nil {
			return n.EntityName()
		}
	}

	return v.VariableTemplateId
}

func (v *TenantVariable) EntitySpaceID() string {
	return v.SpaceId
}

func (v *TenantVariable) AddToEntityMap(m EntityMap) {
	v.entityMap = m

	m[v.EntityID()] = v
}

func (v *TenantVariable) DecryptedValue(key []byte) (string, error) {
	return v.Value.DecryptedValue(key)
}

func IsTenantVariableFilename(f string) bool {
	return strings.HasPrefix(f, "TenantVariables-") && strings.HasSuffix(f, ".json")
}

func ReadTenantVariable(filename string) (*TenantVariable, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret TenantVariable
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	return &ret, nil
}
