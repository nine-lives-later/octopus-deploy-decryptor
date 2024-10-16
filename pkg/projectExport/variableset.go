package projectExport

import (
	"encoding/json"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"os"
	"strings"
)

type VariableSet struct {
	Id        string      `json:"Id"`
	OwnerId   string      `json:"OwnerId"`
	OwnerType string      `json:"OwnerType"`
	SpaceId   string      `json:"SpaceId"`
	Variables []*Variable `json:"Variables"`

	// allow for resolving for variable name
	entityMap EntityMap
}

func (v *VariableSet) EntityID() string {
	return v.Id
}

func (v *VariableSet) EntityName() string {
	// resolve the owner
	if v.entityMap != nil {
		if n := v.entityMap[v.OwnerId]; n != nil {
			return n.EntityName()
		}
	}
	return v.OwnerId
}

func (v *VariableSet) EntitySpaceID() string {
	return v.SpaceId
}

func (v *VariableSet) AddToEntityMap(m EntityMap) {
	v.entityMap = m

	m[v.EntityID()] = v

	for _, vv := range v.Variables {
		vv.AddToEntityMap(m)
	}
}

type Variable struct {
	Id    string          `json:"Id"`
	Name  string          `json:"Name"`
	Type  string          `json:"Type"`
	Value string          `json:"Value"`
	Scope json.RawMessage `json:"Scope"`

	// injected by `ReadVariableSet()`
	SpaceId string `json:"-"`
}

var (
	_ Entity = &VariableSet{} // ensure type compatibility
	_ Entity = &Variable{}    // ensure type compatibility
)

func (v *Variable) EntityID() string {
	return v.Id
}

func (v *Variable) EntityName() string {
	return v.Name
}

func (v *Variable) EntitySpaceID() string {
	return v.SpaceId
}

func (v *Variable) AddToEntityMap(m EntityMap) {
	m[v.EntityID()] = v
}

func (v *Variable) DecryptedValue(key []byte) (string, error) {
	if v.Type == "Sensitive" {
		return decryptor.DecryptString(key, v.Value)
	}

	return v.Value, nil
}

func IsVariableSetFilename(f string) bool {
	return strings.HasPrefix(f, "variableset-") && strings.HasSuffix(f, ".json")
}

func ReadVariableSet(filename string) (*VariableSet, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%v': %w", filename, err)
	}

	var ret VariableSet
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%v': %w", filename, err)
	}

	// inject space into variables
	for _, v := range ret.Variables {
		v.SpaceId = ret.SpaceId
	}

	return &ret, nil
}
