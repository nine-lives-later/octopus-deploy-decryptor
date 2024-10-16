package projectExport

import (
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadVariableSet(t *testing.T) {
	vv, err := ReadVariableSet("../../Octopus-Export-Example/variableset-Projects-441-D4C63A85E7894C5D8C20D9297FEA1A43.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vv)
	assert.True(t, len(vv.Variables) > 0)

	testEntity(t, vv)

	// decrypt all variables
	key, err := decryptor.KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range vv.Variables {
		val, err := v.DecryptedValue(key)
		if err != nil {
			t.Fatal(err)
		}

		testEntity(t, vv)

		t.Logf("%v = %v", v.Name, val)
	}
}
