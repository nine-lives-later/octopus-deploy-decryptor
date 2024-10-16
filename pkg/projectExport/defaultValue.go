package projectExport

import (
	"encoding/json"
	"fmt"
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
)

type DefaultValue struct {
	Value         string
	SenstiveValue string
	Unknown       json.RawMessage
}

var (
	_ json.Unmarshaler = &DefaultValue{} // ensure interface compatibility
)

func (v *DefaultValue) UnmarshalJSON(data []byte) error {
	var val any
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}

	switch vv := val.(type) {
	case string:
		v.Value = vv
	case map[string]interface{}:
		if ss := vv["SensitiveValue"]; ss != nil {
			v.SenstiveValue = ss.(string)
		}
	default:
		return fmt.Errorf("unexpected type '%T' for variable value", val)
	}

	return nil
}

func (v *DefaultValue) DecryptedValue(key []byte) (string, error) {
	if v.SenstiveValue != "" {
		return decryptor.DecryptString(key, v.SenstiveValue)
	}
	if len(v.Unknown) > 0 {
		return "", fmt.Errorf("unknown type for variable value: '%v'", string(v.Unknown))
	}
	return v.Value, nil
}
