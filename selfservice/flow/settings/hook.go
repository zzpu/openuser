package settings

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/identity"
	"github.com/zzpu/openuser/selfservice/flow"
	"github.com/zzpu/openuser/x"
)

type (
	PostHookPrePersistExecutor interface {
		ExecuteSettingsPrePersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error
	}
	PostHookPrePersistExecutorFunc func(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error

	PostHookPostPersistExecutor interface {
		ExecuteSettingsPostPersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error
	}
	PostHookPostPersistExecutorFunc func(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error

	HooksProvider interface {
		PostSettingsPrePersistHooks(settingsType string) []PostHookPrePersistExecutor
		PostSettingsPostPersistHooks(settingsType string) []PostHookPostPersistExecutor
	}
)

func (f PostHookPrePersistExecutorFunc) ExecuteSettingsPrePersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error {
	return f(w, r, a, s)
}

func (f PostHookPostPersistExecutorFunc) ExecuteSettingsPostPersistHook(w http.ResponseWriter, r *http.Request, a *Flow, s *identity.Identity) error {
	return f(w, r, a, s)
}

type (
	executorDependencies interface {
		identity.ManagementProvider
		identity.ValidationProvider
		HooksProvider
		x.LoggingProvider
		FlowPersistenceProvider
		x.WriterProvider
	}
	HookExecutor struct {
		d executorDependencies
		c configuration.Provider
	}
	HookExecutorProvider interface {
		SettingsHookExecutor() *HookExecutor
	}
)

func PostHookPostPersistExecutorNames(e []PostHookPostPersistExecutor) []string {
	names := make([]string, len(e))
	for k, ee := range e {
		names[k] = fmt.Sprintf("%T", ee)
	}
	return names
}

func PostHookPrePersistExecutorNames(e []PostHookPrePersistExecutor) []string {
	names := make([]string, len(e))
	for k, ee := range e {
		names[k] = fmt.Sprintf("%T", ee)
	}
	return names
}

func NewHookExecutor(
	d executorDependencies,
	c configuration.Provider,
) *HookExecutor {
	return &HookExecutor{
		d: d,
		c: c,
	}
}

type PostSettingsHookOption func(o *postSettingsHookOptions)

type postSettingsHookOptions struct {
	cb func(ctxUpdate *UpdateContext) error
}

func WithCallback(cb func(ctxUpdate *UpdateContext) error) func(o *postSettingsHookOptions) {
	return func(o *postSettingsHookOptions) {
		o.cb = cb
	}
}

func (e *HookExecutor) PostSettingsHook(w http.ResponseWriter, r *http.Request, settingsType string, ctxUpdate *UpdateContext, i *identity.Identity, opts ...PostSettingsHookOption) error {
	e.d.Logger().
		WithRequest(r).
		WithField("identity_id", i.ID).
		WithField("flow_method", settingsType).
		Debug("Running PostSettingsPrePersistHooks.")

	config := new(postSettingsHookOptions)
	for _, f := range opts {
		f(config)
	}

	for k, executor := range e.d.PostSettingsPrePersistHooks(settingsType) {
		logFields := logrus.Fields{
			"executor":          fmt.Sprintf("%T", executor),
			"executor_position": k,
			"executors":         PostHookPrePersistExecutorNames(e.d.PostSettingsPrePersistHooks(settingsType)),
			"identity_id":       i.ID,
			"flow_method":       settingsType,
		}

		if err := executor.ExecuteSettingsPrePersistHook(w, r, ctxUpdate.Flow, i); err != nil {
			if errors.Is(err, ErrHookAbortRequest) {
				e.d.Logger().WithRequest(r).WithFields(logFields).
					Debug("A ExecuteSettingsPrePersistHook hook aborted early.")
				return nil
			}
			return err
		}

		e.d.Logger().WithRequest(r).WithFields(logFields).Debug("ExecuteSettingsPrePersistHook completed successfully.")
	}

	options := []identity.ManagerOption{identity.ManagerExposeValidationErrorsForInternalTypeAssertion}
	ttl := e.c.SelfServiceFlowSettingsPrivilegedSessionMaxAge()
	if ctxUpdate.Session.AuthenticatedAt.Add(ttl).After(time.Now()) {
		options = append(options, identity.ManagerAllowWriteProtectedTraits)
	}

	if err := e.d.IdentityManager().Update(r.Context(), i, options...); err != nil {
		if errors.Is(err, identity.ErrProtectedFieldModified) {
			e.d.Logger().WithError(err).Debug("Modifying protected field requires re-authentication.")
			return errors.WithStack(NewFlowNeedsReAuth())
		}
		return err
	}
	e.d.Audit().
		WithRequest(r).
		WithField("identity_id", i.ID).
		Debug("An identity's settings have been updated.")

	ctxUpdate.Session.Identity = i
	ctxUpdate.Flow.State = StateSuccess
	if config.cb != nil {
		if err := config.cb(ctxUpdate); err != nil {
			return err
		}
	}

	if method, ok := ctxUpdate.Flow.Methods[settingsType]; ok {
		method.Config.ResetMessages()
	}

	if err := e.d.SettingsFlowPersister().UpdateSettingsFlow(r.Context(), ctxUpdate.Flow); err != nil {
		return err
	}

	for k, executor := range e.d.PostSettingsPostPersistHooks(settingsType) {
		if err := executor.ExecuteSettingsPostPersistHook(w, r, ctxUpdate.Flow, i); err != nil {
			if errors.Is(err, ErrHookAbortRequest) {
				e.d.Logger().
					WithRequest(r).
					WithField("executor", fmt.Sprintf("%T", executor)).
					WithField("executor_position", k).
					WithField("executors", PostHookPostPersistExecutorNames(e.d.PostSettingsPostPersistHooks(settingsType))).
					WithField("identity_id", i.ID).
					WithField("flow_method", settingsType).
					Debug("A ExecuteSettingsPostPersistHook hook aborted early.")
				return nil
			}
			return err
		}

		e.d.Logger().WithRequest(r).
			WithField("executor", fmt.Sprintf("%T", executor)).
			WithField("executor_position", k).
			WithField("executors", PostHookPostPersistExecutorNames(e.d.PostSettingsPostPersistHooks(settingsType))).
			WithField("identity_id", i.ID).
			WithField("flow_method", settingsType).
			Debug("ExecuteSettingsPostPersistHook completed successfully.")
	}

	e.d.Logger().
		WithRequest(r).
		WithField("identity_id", i.ID).
		WithField("flow_method", settingsType).
		Debug("Completed all PostSettingsPrePersistHooks and PostSettingsPostPersistHooks.")

	if ctxUpdate.Flow.Type == flow.TypeAPI {
		updatedFlow, err := e.d.SettingsFlowPersister().GetSettingsFlow(r.Context(), ctxUpdate.Flow.ID)
		if err != nil {
			return err
		}

		e.d.Writer().Write(w, r, &APIFlowResponse{Flow: updatedFlow, Identity: i})
		return nil
	}

	return x.SecureContentNegotiationRedirection(w, r, ctxUpdate.Session.Declassify(), ctxUpdate.Flow.RequestURL, e.d.Writer(), e.c,
		x.SecureRedirectOverrideDefaultReturnTo(
			e.c.SelfServiceFlowSettingsReturnTo(settingsType,
				ctxUpdate.Flow.AppendTo(e.c.SelfServiceFlowSettingsUI()))))
}
