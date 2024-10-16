package projectExport

import (
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadTenantVariables(t *testing.T) {
	vv, err := ReadTenantVariable("../../Octopus-Export-Example/TenantVariables-762-D4C63A85E7894C5D8C20D9297FEA1A43.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vv)
	assert.NotNil(t, vv.Value)

	testEntity(t, vv)

	// decrypt all variables
	key, err := decryptor.KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	val, err := vv.DecryptedValue(key)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "sensitive_custom", val)
}
