package variableset

import (
	"github.com/nine-lives-later/octopus-deploy-decryptor/pkg/decryptor"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	examplePassword = "pass!w0rd"
)

func TestReadVariables(t *testing.T) {
	vv, err := ReadVariables("../../Octopus-Export-Example/variableset-Projects-441-D4C63A85E7894C5D8C20D9297FEA1A43.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vv)
	assert.True(t, len(vv) > 0)

	// decrypt all variables
	key, err := decryptor.KeyFromPassword(examplePassword)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range vv {
		val, err := v.DecryptedValue(key)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v = %v", v.Name, val)
	}
}
