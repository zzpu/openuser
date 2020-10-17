package password

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/ory/x/decoderx"

	"github.com/zzpu/openuser/continuity"
	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/hash"
	"github.com/zzpu/openuser/identity"
	"github.com/zzpu/openuser/selfservice/errorx"
	"github.com/zzpu/openuser/selfservice/flow/login"
	"github.com/zzpu/openuser/selfservice/flow/registration"
	"github.com/zzpu/openuser/selfservice/flow/settings"
	"github.com/zzpu/openuser/session"
	"github.com/zzpu/openuser/x"
)

var _ login.Strategy = new(Strategy)
var _ registration.Strategy = new(Strategy)
var _ identity.ActiveCredentialsCounter = new(Strategy)

type registrationStrategyDependencies interface {
	x.LoggingProvider
	x.WriterProvider
	x.CSRFTokenGeneratorProvider
	x.CSRFProvider

	continuity.ManagementProvider

	errorx.ManagementProvider
	ValidationProvider
	hash.HashProvider

	registration.HandlerProvider
	registration.HooksProvider
	registration.ErrorHandlerProvider
	registration.HookExecutorProvider
	registration.FlowPersistenceProvider

	login.HooksProvider
	login.ErrorHandlerProvider
	login.HookExecutorProvider
	login.FlowPersistenceProvider
	login.HandlerProvider

	settings.FlowPersistenceProvider
	settings.HookExecutorProvider
	settings.HooksProvider
	settings.ErrorHandlerProvider

	identity.PrivilegedPoolProvider
	identity.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
}

type Strategy struct {
	c  configuration.Provider
	d  registrationStrategyDependencies
	v  *validator.Validate
	hd *decoderx.HTTP
}

func (s *Strategy) CountActiveCredentials(cc map[identity.CredentialsType]identity.Credentials) (count int, err error) {
	for _, c := range cc {
		if c.Type == s.ID() && len(c.Config) > 0 {
			var conf CredentialsConfig
			if err = json.Unmarshal(c.Config, &conf); err != nil {
				return 0, errors.WithStack(err)
			}

			if len(c.Identifiers) > 0 && len(c.Identifiers[0]) > 0 &&
				strings.HasPrefix(conf.HashedPassword, "$argon2id$") {
				count++
			}
		}
	}
	return
}

func NewStrategy(
	d registrationStrategyDependencies,
	c configuration.Provider,
) *Strategy {
	return &Strategy{
		c:  c,
		d:  d,
		v:  validator.New(),
		hd: decoderx.NewHTTP(),
	}
}

func (s *Strategy) ID() identity.CredentialsType {
	return identity.CredentialsTypePassword
}
