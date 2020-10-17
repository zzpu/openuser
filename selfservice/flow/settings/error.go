package settings

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/x/urlx"

	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/identity"
	"github.com/zzpu/openuser/selfservice/errorx"
	"github.com/zzpu/openuser/selfservice/flow"
	"github.com/zzpu/openuser/selfservice/flow/login"
	"github.com/zzpu/openuser/text"
	"github.com/zzpu/openuser/x"
)

var (
	ErrHookAbortRequest = errors.New("aborted settings hook execution")
)

type (
	errorHandlerDependencies interface {
		errorx.ManagementProvider
		x.WriterProvider
		x.LoggingProvider

		HandlerProvider
		FlowPersistenceProvider
	}

	ErrorHandlerProvider interface{ SettingsFlowErrorHandler() *ErrorHandler }

	ErrorHandler struct {
		d errorHandlerDependencies
		c configuration.Provider
	}

	FlowExpiredError struct {
		*herodot.DefaultError
		ago time.Duration
	}

	FlowNeedsReAuth struct {
		*herodot.DefaultError
	}
)

func NewFlowNeedsReAuth() *FlowNeedsReAuth {
	return &FlowNeedsReAuth{DefaultError: herodot.ErrForbidden.
		WithReasonf("The login session is too old and thus not allowed to update these fields. Please re-authenticate.")}
}

func NewFlowExpiredError(at time.Time) *FlowExpiredError {
	ago := time.Since(at)
	return &FlowExpiredError{
		ago: ago,
		DefaultError: herodot.ErrBadRequest.
			WithError("settings flow expired").
			WithReasonf(`The settings flow has expired. Please restart the flow.`).
			WithReasonf("The settings flow expired %.2f minutes ago, please try again.", ago.Minutes()),
	}
}

func NewErrorHandler(d errorHandlerDependencies, c configuration.Provider) *ErrorHandler {
	return &ErrorHandler{
		d: d,
		c: c,
	}
}

func (s *ErrorHandler) reauthenticate(
	w http.ResponseWriter,
	r *http.Request,
	f *Flow,
	err error,
) {
	if f.Type == flow.TypeAPI {
		s.d.Writer().WriteError(w, r, err)
		return
	}

	returnTo := urlx.CopyWithQuery(urlx.AppendPaths(s.c.SelfPublicURL(), r.URL.Path), r.URL.Query())
	http.Redirect(w, r, urlx.AppendPaths(urlx.CopyWithQuery(s.c.SelfPublicURL(),
		url.Values{
			"refresh":   {"true"},
			"return_to": {returnTo.String()},
		}), login.RouteInitBrowserFlow).String(), http.StatusFound)
}

func (s *ErrorHandler) WriteFlowError(
	w http.ResponseWriter,
	r *http.Request,
	method string,
	f *Flow,
	id *identity.Identity,
	err error,
) {
	s.d.Audit().
		WithError(err).
		WithRequest(r).
		WithField("settings_flow", f).
		Info("Encountered self-service settings error.")

	if f == nil {
		s.forward(w, r, f, err)
		return
	}

	if e := new(FlowExpiredError); errors.As(err, &e) {
		// create new flow because the old one is not valid
		a, err := s.d.SettingsHandler().NewFlow(w, r, id, f.Type)
		if err != nil {
			// failed to create a new session and redirect to it, handle that error as a new one
			s.WriteFlowError(w, r, method, f, id, err)
			return
		}

		a.Messages.Add(text.NewErrorValidationSettingsFlowExpired(e.ago))
		if err := s.d.SettingsFlowPersister().UpdateSettingsFlow(r.Context(), a); err != nil {
			s.forward(w, r, a, err)
			return
		}

		if f.Type == flow.TypeAPI {
			http.Redirect(w, r, urlx.CopyWithQuery(urlx.AppendPaths(s.c.SelfPublicURL(),
				RouteGetFlow), url.Values{"id": {a.ID.String()}}).String(), http.StatusFound)
		} else {
			http.Redirect(w, r, a.AppendTo(s.c.SelfServiceFlowSettingsUI()).String(), http.StatusFound)
		}
		return
	}

	if e := new(FlowNeedsReAuth); errors.As(err, &e) {
		s.reauthenticate(w, r, f, err)
		return
	}

	if _, ok := f.Methods[method]; !ok {
		s.forward(w, r, f, errors.WithStack(herodot.ErrInternalServerError.
			WithErrorf(`Expected settings method "%s" to exist in flow. This is a bug in the code and should be reported on GitHub.`, method)))
		return
	}

	if err := f.Methods[method].Config.ParseError(err); err != nil {
		s.forward(w, r, f, err)
		return
	}

	if err := s.d.SettingsFlowPersister().UpdateSettingsFlowMethod(r.Context(), f.ID, method, f.Methods[method]); err != nil {
		s.forward(w, r, f, err)
		return
	}

	if f.Type == flow.TypeBrowser {
		http.Redirect(w, r, f.AppendTo(s.c.SelfServiceFlowSettingsUI()).String(), http.StatusFound)
		return
	}

	updatedFlow, innerErr := s.d.SettingsFlowPersister().GetSettingsFlow(r.Context(), f.ID)
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
