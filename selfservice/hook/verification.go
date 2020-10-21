package hook

import (
	"net/http"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/flow/registration"
	"github.com/zzpu/ums/selfservice/flow/settings"
	"github.com/zzpu/ums/selfservice/strategy/link"
	"github.com/zzpu/ums/session"
)

var _ registration.PostHookPostPersistExecutor = new(Verifier)
var _ settings.PostHookPostPersistExecutor = new(Verifier)

type (
	verifierDependencies interface {
		link.SenderProvider
		link.VerificationTokenPersistenceProvider
	}
	Verifier struct {
		r verifierDependencies
		c configuration.Provider
	}
)

func NewVerifier(r verifierDependencies, c configuration.Provider) *Verifier {
	return &Verifier{r: r, c: c}
}

func (e *Verifier) ExecutePostRegistrationPostPersistHook(_ http.ResponseWriter, r *http.Request, _ *registration.Flow, s *session.Session) error {
	return e.do(r, s.Identity)
}

func (e *Verifier) ExecuteSettingsPostPersistHook(w http.ResponseWriter, r *http.Request, a *settings.Flow, i *identity.Identity) error {
	return e.do(r, i)
}

func (e *Verifier) do(r *http.Request, i *identity.Identity) error {
	// Ths is called after the identity has been created so we can safely assume that all addresses are available
	// already.

	for k := range i.VerifiableAddresses {
		address := &i.VerifiableAddresses[k]
		if address.Verified {
			continue
		}

		token := link.NewVerificationToken(address, e.c.SelfServiceFlowVerificationRequestLifespan())
		if err := e.r.VerificationTokenPersister().CreateVerificationToken(r.Context(), token); err != nil {
			return err
		}

		if err := e.r.LinkSender().SendVerificationTokenTo(r.Context(), address, token); err != nil {
			return err
		}
	}

	return nil
}
