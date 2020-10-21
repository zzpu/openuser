package identity

import (
	"testing"

	"github.com/zzpu/ums/driver/configuration"

	"github.com/stretchr/testify/assert"
)

func TestNewIdentity(t *testing.T) {
	i := NewIdentity(configuration.DefaultIdentityTraitsSchemaID)
	assert.NotEmpty(t, i.ID)
	// assert.NotEmpty(t, i.Metadata)
	assert.NotEmpty(t, i.Traits)
	assert.NotNil(t, i.Credentials)
}
