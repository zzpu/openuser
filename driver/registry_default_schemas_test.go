package driver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/viper"
	"github.com/ory/x/urlx"
	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/internal"
	"github.com/zzpu/openuser/schema"
)

func TestRegistryDefault_IdentityTraitsSchemas(t *testing.T) {
	_, reg := internal.NewFastRegistryWithMocks(t)
	defaultSchema := schema.Schema{
		ID:     "default",
		URL:    urlx.ParseOrPanic("file://default.schema.json"),
		RawURL: "file://default.schema.json",
	}
	altSchema := schema.Schema{
		ID:     "alt",
		URL:    urlx.ParseOrPanic("file://other.schema.json"),
		RawURL: "file://other.schema.json",
	}

	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, defaultSchema.RawURL)
	viper.Set(configuration.ViperKeyIdentitySchemas, []configuration.SchemaConfig{{ID: altSchema.ID, URL: altSchema.RawURL}})

	ss := reg.IdentityTraitsSchemas()
	assert.Equal(t, 2, len(ss))
	assert.Contains(t, ss, defaultSchema)
	assert.Contains(t, ss, altSchema)
}
