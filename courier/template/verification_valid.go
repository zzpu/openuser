package template

import (
	"path/filepath"

	"github.com/zzpu/ums/driver/configuration"
)

type (
	VerificationValid struct {
		c configuration.Provider
		m *VerificationValidModel
	}
	VerificationValidModel struct {
		To              string
		VerificationURL string
	}
)

func NewVerificationValid(c configuration.Provider, m *VerificationValidModel) *VerificationValid {
	return &VerificationValid{c: c, m: m}
}

func (t *VerificationValid) EmailRecipient() (string, error) {
	return t.m.To, nil
}

func (t *VerificationValid) EmailSubject() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "verification/valid/email.subject.gotmpl"), t.m)
}

func (t *VerificationValid) EmailBody() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "verification/valid/email.body.gotmpl"), t.m)
}
