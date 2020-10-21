package driver

import (
	"github.com/ory/x/tracing"
	"github.com/zzpu/ums/metrics/prometheus"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"

	"github.com/ory/x/logrusx"

	"github.com/zzpu/ums/continuity"
	"github.com/zzpu/ums/courier"
	"github.com/zzpu/ums/hash"
	"github.com/zzpu/ums/schema"
	"github.com/zzpu/ums/selfservice/flow/recovery"
	"github.com/zzpu/ums/selfservice/flow/settings"
	"github.com/zzpu/ums/selfservice/flow/verification"
	"github.com/zzpu/ums/selfservice/strategy/link"

	"github.com/ory/x/healthx"

	"github.com/zzpu/ums/persistence"
	"github.com/zzpu/ums/selfservice/flow/login"
	"github.com/zzpu/ums/selfservice/flow/logout"
	"github.com/zzpu/ums/selfservice/flow/registration"

	"github.com/zzpu/ums/x"

	"github.com/ory/x/dbal"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/errorx"
	password2 "github.com/zzpu/ums/selfservice/strategy/password"
	"github.com/zzpu/ums/session"
)

type Registry interface {
	dbal.Driver

	Init() error

	WithConfig(c configuration.Provider) Registry
	WithLogger(l *logrusx.Logger) Registry

	BuildVersion() string
	BuildDate() string
	BuildHash() string
	WithBuildInfo(version, hash, date string) Registry

	WithCSRFHandler(c x.CSRFHandler)
	WithCSRFTokenGenerator(cg x.CSRFToken)

	HealthHandler() *healthx.Handler
	CookieManager() sessions.Store
	ContinuityCookieManager() sessions.Store

	RegisterRoutes(public *x.RouterPublic, admin *x.RouterAdmin)
	RegisterPublicRoutes(public *x.RouterPublic)
	RegisterAdminRoutes(admin *x.RouterAdmin)
	PrometheusManager() *prometheus.MetricsManager
	Tracer() *tracing.Tracer

	x.CSRFProvider
	x.WriterProvider
	x.LoggingProvider

	continuity.ManagementProvider
	continuity.PersistenceProvider

	courier.Provider

	persistence.Provider

	errorx.ManagementProvider
	errorx.HandlerProvider
	errorx.PersistenceProvider

	hash.HashProvider

	identity.HandlerProvider
	identity.ValidationProvider
	identity.PoolProvider
	identity.PrivilegedPoolProvider
	identity.ManagementProvider
	identity.ActiveCredentialsCounterStrategyProvider

	schema.HandlerProvider

	password2.ValidationProvider

	session.HandlerProvider
	session.ManagementProvider
	session.PersistenceProvider

	settings.HandlerProvider
	settings.ErrorHandlerProvider
	settings.FlowPersistenceProvider
	settings.StrategyProvider

	login.FlowPersistenceProvider
	login.ErrorHandlerProvider
	login.HooksProvider
	login.HookExecutorProvider
	login.HandlerProvider
	login.StrategyProvider

	logout.HandlerProvider

	registration.FlowPersistenceProvider
	registration.ErrorHandlerProvider
	registration.HooksProvider
	registration.HookExecutorProvider
	registration.HandlerProvider
	registration.StrategyProvider

	verification.FlowPersistenceProvider
	verification.ErrorHandlerProvider
	verification.HandlerProvider
	verification.StrategyProvider

	link.SenderProvider
	link.VerificationTokenPersistenceProvider
	link.RecoveryTokenPersistenceProvider

	recovery.FlowPersistenceProvider
	recovery.ErrorHandlerProvider
	recovery.HandlerProvider
	recovery.StrategyProvider

	x.CSRFTokenGeneratorProvider
}

func NewRegistry(c configuration.Provider) (Registry, error) {
	dsn := c.DSN()
	driver, err := dbal.GetDriverFor(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	return registry, nil
}
