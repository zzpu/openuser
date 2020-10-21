package link

import (
	"github.com/ory/x/decoderx"

	"github.com/zzpu/ums/courier"
	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/schema"
	"github.com/zzpu/ums/selfservice/errorx"
	"github.com/zzpu/ums/selfservice/flow/recovery"
	"github.com/zzpu/ums/selfservice/flow/settings"
	"github.com/zzpu/ums/selfservice/flow/verification"
	"github.com/zzpu/ums/selfservice/form"
	"github.com/zzpu/ums/session"
	"github.com/zzpu/ums/x"
)

var _ recovery.Strategy = new(Strategy)
var _ recovery.AdminHandler = new(Strategy)
var _ recovery.PublicHandler = new(Strategy)

var _ verification.Strategy = new(Strategy)
var _ verification.AdminHandler = new(Strategy)
var _ verification.PublicHandler = new(Strategy)

type (
	// FlowMethod contains the configuration for this selfservice strategy.
	FlowMethod struct {
		*form.HTMLForm
	}

	strategyDependencies interface {
		x.CSRFProvider
		x.CSRFTokenGeneratorProvider
		x.WriterProvider
		x.LoggingProvider

		session.HandlerProvider
		session.ManagementProvider
		settings.HandlerProvider
		settings.FlowPersistenceProvider

		identity.ValidationProvider
		identity.ManagementProvider
		identity.PoolProvider
		identity.PrivilegedPoolProvider

		courier.Provider

		errorx.ManagementProvider

		recovery.ErrorHandlerProvider
		recovery.FlowPersistenceProvider
		recovery.StrategyProvider

		verification.ErrorHandlerProvider
		verification.FlowPersistenceProvider
		verification.StrategyProvider

		RecoveryTokenPersistenceProvider
		VerificationTokenPersistenceProvider
		SenderProvider

		IdentityTraitsSchemas() schema.Schemas
	}

	Strategy struct {
		c  configuration.Provider
		d  strategyDependencies
		dx *decoderx.HTTP
	}
)

func NewStrategy(d strategyDependencies, c configuration.Provider) *Strategy {
	return &Strategy{c: c, d: d, dx: decoderx.NewHTTP()}
}
