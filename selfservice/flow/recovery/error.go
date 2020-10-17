package recovery

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/x/sqlxx"

	"github.com/ory/herodot"
	"github.com/ory/x/urlx"

	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/selfservice/errorx"
	"github.com/zzpu/openuser/selfservice/flow"
	"github.com/zzpu/openuser/text"
	"github.com/zzpu/openuser/x"
)

var (
	ErrAlreadyLoggedIn = herodot.ErrBadRequest.WithReason("A valid session was detected and thus recovery is not possible.")
)

type FlowExpiredError struct {
	*herodot.DefaultError
	ago time.Duration
}

func NewFlowExpiredError(at time.Time) *FlowExpiredError {
	ago := time.Since(at)
	return &FlowExpiredError{
		ago: ago,
		DefaultError: herodot.ErrBadRequest.
			WithError("recovery flow expired").
			WithReasonf(`The recovery flow has expired. Please restart the flow.`).
			WithReasonf("The recovery flow expired %.2f minutes ago, please try again.", ago.Minutes()),
	}
}

type (
	errorHandlerDependencies interface {
		errorx.ManagementProvider
		x.WriterProvider
		x.LoggingProvider
		x.CSRFTokenGeneratorProvider
		StrategyProvider

		FlowPersistenceProvider
	}

	ErrorHandlerProvider interface {
		RecoveryFlowErrorHandler() *ErrorHandler
	}

	ErrorHandler struct {
		d errorHandlerDependencies
		c configuration.Provider
	}
)

func NewErrorHandler(d errorHandlerDependencies, c configuration.Provider) *ErrorHandler {
	return &ErrorHandler{
		d: d,
		c: c,
	}
}

func (s *ErrorHandler) WriteFlowError(
	w http.ResponseWriter,
	r *http.Request,
	methodName string,
	f *Flow,
	err error,
) {
	s.d.Audit().
		WithError(err).
		WithRequest(r).
		WithField("recovery_flow", f).
		Info("Encountered self-service recovery error.")

	if f == nil {
		s.forward(w, r, nil, err)
		return
	}

	if e := new(FlowExpiredError); errors.As(err, &e) {
		// create new flow because the old one is not valid
		a, err := NewFlow(s.c.SelfServiceFlowRecoveryRequestLifespan(), s.d.GenerateCSRFToken(r), r, s.d.RecoveryStrategies(), f.Type)
		if err != nil {
			// failed to create a new session and redirect to it, handle that error as a new one
			s.WriteFlowError(w, r, methodName, f, err)
			return
		}

		a.Messages.Add(text.NewErrorValidationRecoveryFlowExpired(e.ago))
		if err := s.d.RecoveryFlowPersister().CreateRecoveryFlow(r.Context(), a); err != nil {
			s.forward(w, r, a, err)
			return
		}

		if f.Type == flow.TypeAPI {
			http.Redirect(w, r, urlx.CopyWithQuery(urlx.AppendPaths(s.c.SelfPublicURL(),
				RouteGetFlow), url.Values{"id": {a.ID.String()}}).String(), http.StatusFound)
		} else {
			http.Redirect(w, r, a.AppendTo(s.c.SelfServiceFlowRecoveryUI()).String(), http.StatusFound)
		}
		return
	}

	method, ok := f.Methods[methodName]
	if !ok {
		s.forward(w, r, f, errors.WithStack(herodot.ErrInternalServerError.
			WithErrorf(`Expected recovery method "%s" to exist in flow. This is a bug in the code and should be reported on GitHub.`, methodName)))
		return
	}

	if err := method.Config.ParseError(err); err != nil {
		s.forward(w, r, f, err)
		return
	}

	f.Active = sqlxx.NullString(methodName)
	if err := s.d.RecoveryFlowPersister().UpdateRecoveryFlow(r.Context(), f); err != nil {
		s.forward(w, r, f, err)
		return
	}

	if f.Type == flow.TypeBrowser {
		http.Redirect(w, r, f.AppendTo(s.c.SelfServiceFlowRecoveryUI()).String(), http.StatusFound)
		return
	}

	updatedFlow, innerErr := s.d.RecoveryFlowPersister().GetRecoveryFlow(r.Context(), f.ID)
	if innerErr != nil {
		s.forward(w, r, updatedFlow, innerErr)
	}

	s.d.Writer().WriteCode(w, r, x.RecoverStatusCode(err, http.StatusBadRequest), updatedFlow)
}

func (s *ErrorHandler) forward(w http.ResponseWriter, r *http.Request, rr *Flow, err error) {
	if rr == nil {
		if x.IsJSONRequest(r) {
			s.d.Writer().WriteError(w, r, err)
			return
		}
		s.d.SelfServiceErrorManager().Forward(r.Context(), w, r, err)
		return
	}

	if rr.Type == flow.TypeAPI {
		s.d.Writer().WriteErrorCode(w, r, x.RecoverStatusCode(err, http.StatusBadRequest), err)
	} else {
		s.d.SelfServiceErrorManager().Forward(r.Context(), w, r, err)
	}
}
