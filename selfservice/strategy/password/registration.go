package password

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"

	"github.com/ory/x/errorsx"

	_ "github.com/ory/jsonschema/v3/fileloader"
	_ "github.com/ory/jsonschema/v3/httploader"
	"github.com/ory/x/decoderx"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/selfservice/flow"
	"github.com/zzpu/ums/session"

	"github.com/ory/herodot"
	"github.com/ory/x/urlx"

	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/schema"
	"github.com/zzpu/ums/selfservice/flow/registration"
	"github.com/zzpu/ums/selfservice/form"
	"github.com/zzpu/ums/x"
)

const (
	RouteRegistration = "/self-service/registration/methods/password"
)

type RegistrationFormPayload struct {
	Password  string          `json:"password"`
	Traits    json.RawMessage `json:"traits"`
	CSRFToken string          `json:"csrf_token"`
}

func (s *Strategy) RegisterRegistrationRoutes(public *x.RouterPublic) {
	s.d.CSRFHandler().ExemptPath(RouteRegistration)
	public.POST(RouteRegistration, s.d.SessionHandler().IsNotAuthenticated(s.handleRegistration, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handler := session.RedirectOnAuthenticated(s.c)
		if x.IsJSONRequest(r) {
			handler = session.RespondWithJSONErrorOnAuthenticated(s.d.Writer(), registration.ErrAlreadyLoggedIn)
		}

		handler(w, r, ps)
	}))
}

func (s *Strategy) handleRegistrationError(w http.ResponseWriter, r *http.Request, rr *registration.Flow, p *RegistrationFormPayload, err error) {
	if rr != nil {
		if method, ok := rr.Methods[identity.CredentialsTypePassword]; ok {
			method.Config.Reset()

			if p != nil {
				for _, field := range form.NewHTMLFormFromJSON("", p.Traits, "traits").Fields {
					// we only set the value and not the whole field because we want to keep types from the initial form generation
					method.Config.SetValue(field.Name, field.Value)
				}
			}

			method.Config.SetCSRF(s.d.GenerateCSRFToken(r))
			rr.Methods[identity.CredentialsTypePassword] = method
			if errSec := method.Config.SortFields(s.c.DefaultIdentityTraitsSchemaURL().String()); errSec != nil {
				s.d.RegistrationFlowErrorHandler().WriteFlowError(w, r, identity.CredentialsTypePassword, rr, errors.Wrap(err, errSec.Error()))
				return
			}
		}
	}

	s.d.RegistrationFlowErrorHandler().WriteFlowError(w, r, identity.CredentialsTypePassword, rr, err)
}

func (s *Strategy) decode(p *RegistrationFormPayload, r *http.Request) error {
	raw, err := sjson.SetBytes(x.MustPkgerRead(pkger.Open("/selfservice/strategy/password/.schema/registration.schema.json")),
		"properties.traits.$ref", s.c.DefaultIdentityTraitsSchemaURL().String()+"#/properties/traits")
	if err != nil {
		return errors.WithStack(err)
	}

	compiler, err := decoderx.HTTPRawJSONSchemaCompiler(raw)
	if err != nil {
		return errors.WithStack(err)
	}

	return s.hd.Decode(r, p, compiler, decoderx.HTTPDecoderSetValidatePayloads(false), decoderx.HTTPDecoderJSONFollowsFormFormat())
}

// nolint:deadcode,unused
// swagger:parameters completeSelfServiceRegistrationFlowWithPasswordMethod
type completeSelfServiceRegistrationFlowWithPasswordMethodParameters struct {
	// Flow is flow ID.
	//
	// in: query
	Flow string `json:"flow"`

	// in: body
	Payload map[string]interface{}
}

// swagger:route POST /self-service/registration/methods/password public completeSelfServiceRegistrationFlowWithPasswordMethod
//
// Complete Registration Flow with Username/Email Password Method
//
// Use this endpoint to complete a registration flow by sending an identity's traits and password. This endpoint
// behaves differently for API and browser flows.
//
// API flows expect `application/json` to be sent in the body and respond with
//   - HTTP 200 and a application/json body with the created identity success - if the session hook is configured the
//     `session` and `session_token` will also be included;
//   - HTTP 302 redirect to a fresh registration flow if the original flow expired with the appropriate error messages set;
//   - HTTP 400 on form validation errors.
//
// Browser flows expect `application/x-www-form-urlencoded` to be sent in the body and responds with
//   - a HTTP 302 redirect to the post/after registration URL or the `return_to` value if it was set and if the registration succeeded;
//   - a HTTP 302 redirect to the registration UI URL with the flow ID containing the validation errors otherwise.
//
// More information can be found at [ORY Kratos User Login and User Registration Documentation](https://www.ory.sh/docs/next/kratos/self-service/flows/user-login-user-registration).
//
//     Schemes: http, https
//
//     Consumes:
//     - application/json
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: registrationViaApiResponse
//       302: emptyResponse
//       400: registrationFlow
//       500: genericError
func (s *Strategy) handleRegistration(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rid := x.ParseUUID(r.URL.Query().Get("flow"))
	if x.IsZeroUUID(rid) {
		s.handleRegistrationError(w, r, nil, nil, errors.WithStack(herodot.ErrBadRequest.WithReasonf("The flow query parameter is missing.")))
		return
	}

	ar, err := s.d.RegistrationFlowPersister().GetRegistrationFlow(r.Context(), rid)
	if err != nil {
		s.handleRegistrationError(w, r, nil, nil, err)
		return
	}

	if err := ar.Valid(); err != nil {
		s.handleRegistrationError(w, r, ar, nil, err)
		return
	}

	var p RegistrationFormPayload
	if err := s.decode(&p, r); err != nil {
		s.handleRegistrationError(w, r, ar, &p, err)
		return
	}

	if err := flow.VerifyRequest(r, ar.Type, s.d.GenerateCSRFToken, p.CSRFToken); err != nil {
		s.handleRegistrationError(w, r, ar, &p, err)
		return
	}

	if len(p.Password) == 0 {
		s.handleRegistrationError(w, r, ar, &p, schema.NewRequiredError("#/password", "password"))
		return
	}

	if len(p.Traits) == 0 {
		p.Traits = json.RawMessage("{}")
	}

	hpw, err := s.d.Hasher().Generate([]byte(p.Password))
	if err != nil {
		s.handleRegistrationError(w, r, ar, &p, err)
		return
	}

	co, err := json.Marshal(&CredentialsConfig{HashedPassword: string(hpw)})
	if err != nil {
		s.handleRegistrationError(w, r, ar, &p, errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Unable to encode password options to JSON: %s", err)))
		return
	}

	i := identity.NewIdentity(configuration.DefaultIdentityTraitsSchemaID)
	i.Traits = identity.Traits(p.Traits)
	i.SetCredentials(s.ID(), identity.Credentials{Type: s.ID(), Identifiers: []string{}, Config: co})

	if err := s.validateCredentials(i, p.Password); err != nil {
		s.handleRegistrationError(w, r, ar, &p, err)
		return
	}

	if err := s.d.RegistrationExecutor().PostRegistrationHook(w, r, identity.CredentialsTypePassword, ar, i); err != nil {
		s.handleRegistrationError(w, r, ar, &p, err)
		return
	}
}

func (s *Strategy) validateCredentials(i *identity.Identity, pw string) error {
	if err := s.d.IdentityValidator().Validate(i); err != nil {
		return err
	}

	c, ok := i.GetCredentials(identity.CredentialsTypePassword)
	if !ok {
		// This should never happen
		return errors.WithStack(x.PseudoPanic.WithReasonf("identity object did not provide the %s CredentialType unexpectedly", identity.CredentialsTypePassword))
	} else if len(c.Identifiers) == 0 {
		return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("No login identifiers (e.g. email, phone number, username) were set. Contact an administrator, the identity schema is misconfigured."))
	}

	for _, id := range c.Identifiers {
		if err := s.d.PasswordValidator().Validate(id, pw); err != nil {
			if _, ok := errorsx.Cause(err).(*herodot.DefaultError); ok {
				return err
			}
			return schema.NewPasswordPolicyViolationError("#/password", err.Error())
		}
	}

	return nil
}

func (s *Strategy) PopulateRegistrationMethod(r *http.Request, sr *registration.Flow) error {
	action := sr.AppendTo(urlx.AppendPaths(s.c.SelfPublicURL(), RouteRegistration))

	htmlf, err := form.NewHTMLFormFromJSONSchema(action.String(), s.c.DefaultIdentityTraitsSchemaURL().String(), "", nil)
	if err != nil {
		return err
	}

	htmlf.Method = "POST"
	htmlf.SetCSRF(s.d.GenerateCSRFToken(r))
	htmlf.SetField(form.Field{Name: "password", Type: "password", Required: true})

	if err := htmlf.SortFields(s.c.DefaultIdentityTraitsSchemaURL().String()); err != nil {
		return err
	}

	sr.Methods[identity.CredentialsTypePassword] = &registration.FlowMethod{
		Method: identity.CredentialsTypePassword,
		Config: &registration.FlowMethodConfig{FlowMethodConfigurator: &FlowMethod{HTMLForm: htmlf}},
	}

	return nil
}
