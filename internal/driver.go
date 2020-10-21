package internal

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"github.com/ory/x/logrusx"

	"github.com/zzpu/ums/driver"
	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/selfservice/hook"
	"github.com/zzpu/ums/x"
)

func resetConfig() {
	viper.Set(configuration.ViperKeyDSN, nil)

	viper.Set("LOG_LEVEL", "trace")
	viper.Set(configuration.ViperKeyHasherArgon2ConfigMemory, 64)
	viper.Set(configuration.ViperKeyHasherArgon2ConfigIterations, 1)
	viper.Set(configuration.ViperKeyHasherArgon2ConfigParallelism, 1)
	viper.Set(configuration.ViperKeyHasherArgon2ConfigSaltLength, 2)
	viper.Set(configuration.ViperKeyHasherArgon2ConfigKeyLength, 2)
}

func NewConfigurationWithDefaults() *configuration.ViperProvider {
	viper.Reset()
	resetConfig()
	return configuration.NewViperProvider(logrusx.New("", ""), true)
}

// NewFastRegistryWithMocks returns a registry with several mocks and an SQLite in memory database that make testing
// easier and way faster. This suite does not work for e2e or advanced integration tests.
func NewFastRegistryWithMocks(t *testing.T) (*configuration.ViperProvider, *driver.RegistryDefault) {
	conf, reg := NewRegistryDefaultWithDSN(t, "")
	reg.WithCSRFTokenGenerator(x.FakeCSRFTokenGenerator)
	reg.WithCSRFHandler(x.NewFakeCSRFHandler(""))
	reg.WithHooks(map[string]func(configuration.SelfServiceHook) interface{}{
		"err": func(c configuration.SelfServiceHook) interface{} {
			return &hook.Error{Config: c.Config}
		},
	})

	require.NoError(t, reg.Persister().MigrateUp(context.Background()))
	return conf, reg
}

// NewRegistryDefaultWithDSN returns a more standard registry without mocks. Good for e2e and advanced integration testing!
func NewRegistryDefaultWithDSN(t *testing.T, dsn string) (*configuration.ViperProvider, *driver.RegistryDefault) {
	viper.Reset()
	resetConfig()

	viper.Set(configuration.ViperKeyDSN, "sqlite3://"+filepath.Join(os.TempDir(), x.NewUUID().String())+".sql?mode=memory&_fk=true")
	if dsn != "" {
		viper.Set(configuration.ViperKeyDSN, dsn)
	}

	d, err := driver.NewDefaultDriver(logrusx.New("", ""), "test", "test", "test", true)
	require.NoError(t, err)
	return d.Configuration().(*configuration.ViperProvider), d.Registry().(*driver.RegistryDefault)
}
