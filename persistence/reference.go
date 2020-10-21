package persistence

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"

	"github.com/zzpu/ums/continuity"
	"github.com/zzpu/ums/courier"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/errorx"
	"github.com/zzpu/ums/selfservice/flow/login"
	"github.com/zzpu/ums/selfservice/flow/recovery"
	"github.com/zzpu/ums/selfservice/flow/registration"
	"github.com/zzpu/ums/selfservice/flow/settings"
	"github.com/zzpu/ums/selfservice/flow/verification"
	"github.com/zzpu/ums/selfservice/strategy/link"
	"github.com/zzpu/ums/session"
)

type Provider interface {
	Persister() Persister
}

type Persister interface {
	continuity.Persister
	identity.PrivilegedPool
	registration.FlowPersister
	login.FlowPersister
	settings.FlowPersister
	courier.Persister
	session.Persister
	errorx.Persister
	verification.FlowPersister
	recovery.FlowPersister
	link.RecoveryTokenPersister
	link.VerificationTokenPersister

	Close(context.Context) error
	Ping(context.Context) error
	MigrationStatus(c context.Context, b io.Writer) error
	MigrateDown(c context.Context, steps int) error
	MigrateUp(c context.Context) error
	GetConnection(ctx context.Context) *pop.Connection
	Transaction(ctx context.Context, callback func(ctx context.Context, connection *pop.Connection) error) error
}
