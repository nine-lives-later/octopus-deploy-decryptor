package projectExport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testEntity(t *testing.T, e Entity) {
	assert.NotNil(t, e)
	assert.NotEmpty(t, e.EntityID())
	assert.NotEmpty(t, e.EntityName())
	assert.NotEmpty(t, e.EntitySpaceID())
}
