package login

import (
	"net/http"
	"net/url"
	"time"

	"github.com/ory/x/urlx"

	"github.com/zzpu/ums/selfservice/flow"
	"github.com/zzpu/ums/text"

	"github.com/pkg/errors"

	"github.com/ory/herodot"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/errorx"
	"github.com/zzpu/ums/x"
)

var (
	ErrHookAbortFlow   = errors.New("aborted login hook execution")
	ErrAlreadyLoggedIn = herodot.ErrBadRequest.WithReason("A valid session was detected and thus login is not possible. Did you forget to set `?refresh=true`?")
)

type (
	errorHandlerDependencies interface {
		errorx.ManagementProvider
		x.WriterProvider
		x.LoggingProvider

		FlowPersistenceProvider
		HandlerProvider
	}

	ErrorHandlerProvider interface{ LoginFlowErrorHandler() *ErrorHandler }

	ErrorHandler struct {
		d errorHandlerDependencies
		c configuration.Provider
	}

	FlowExpiredError struct {
		*herodot.DefaultError
		ago time.Duration
	}
)

func NewFlowExpiredError(at time.Time) *FlowExpiredError {
	ago := time.Since(at)
	return &FlowExpiredError{
		ago: ago,
		DefaultError: herodot.ErrBadRequest.
			WithError("login flow expired").
			WithReasonf(`The login flow has expired. Please restart the flow.`).
			WithReasonf("The login flow expired %.2f minutes ago, please try again.", ago.Minutes()),
	}
}

func NewFlowErrorHandler(d errorHandlerDependencies, c configuration.Provider) *ErrorHandler {
	return &ErrorHandler{d: d, c: c}
}

func (s *ErrorHandler) WriteFlowError(w http.ResponseWriter, r *http.Request, ct identity.CredentialsType, f *Flow, err error) {
	s.d.Audit().
		WithError(err).
		WithRequest(r).
		WithField("login_flow", f).
		Info("Encountered self-service login error.")

	if f == nil {
		s.forward(w, r, nil, err)
		return
	}

	if e := new(FlowExpiredError); errors.As(err, &e) {
		// create new flow because the old one is not valid
		a, err := s.d.LoginHandler().NewLoginFlow(w, r, f.Type)
		if err != nil {
			// failed to create a new session and redirect to it, handle that error as a new one
			s.WriteFlowError(w, r, ct, f, err)
			return
		}

		a.Messages.Add(text.NewErrorValidationLoginFlowExpired(e.ago))
		if err := s.d.LoginFlowPersister().UpdateLoginFlow(r.Context(), a); err != nil {
			s.forward(w, r, a, err)
			return
		}

		if f.Type == flow.TypeAPI {
			http.Redirect(w, r, urlx.CopyWithQuery(urlx.AppendPaths(s.c.SelfPublicURL(),
				RouteGetFlow), url.Values{"id": {a.ID.String()}}).String(), http.StatusFound)
		} else {
			http.Redirect(w, r, a.AppendTo(s.c.SelfServiceFlowLoginUI()).String(), http.StatusFound)
		}
		return
	}

	method, ok := f.Methods[ct]
	if !ok {
		s.forward(w, r, f, errors.WithStack(herodot.ErrInternalServerError.
			WithErrorf(`Expected login method "%s" to exist in flow. This is a bug in the code and should be reported on GitHub.`, ct)))
		return
	}

	if err := method.Config.ParseError(err); err != nil {
		s.forward(w, r, f, err)
		return
	}

	if err := s.d.LoginFlowPersister().UpdateLoginFlowMethod(r.Context(), f.ID, ct, method); err != nil {
		s.forward(w, r, f, err)
		return
	}

	if f.Type == flow.TypeBrowser {
		http.Redirect(w, r, f.AppendTo(s.c.SelfServiceFlowLoginUI()).String(), http.StatusFound)
		return
	}

	updatedFlow, innerErr := s.d.LoginFlowPersister().GetLoginFlow(r.Context(), f.ID)
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
