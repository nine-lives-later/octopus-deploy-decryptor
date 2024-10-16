package projectExport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadTenant(t *testing.T) {
	vv, err := ReadTenant("../../Octopus-Export-Example/Tenants-21-D4C63A85E7894C5D8C20D9297FEA1A43.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vv)

	testEntity(t, vv)
}
