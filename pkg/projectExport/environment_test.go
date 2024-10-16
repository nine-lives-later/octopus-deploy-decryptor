package projectExport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadEnvironment(t *testing.T) {
	vv, err := ReadEnvironment("../../Octopus-Export-Example/Environments-1-D4C63A85E7894C5D8C20D9297FEA1A43.json")
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, vv)

	testEntity(t, vv)
}
