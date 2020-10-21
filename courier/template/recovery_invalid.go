package template

import (
	"path/filepath"

	"github.com/zzpu/ums/driver/configuration"
)

type (
	RecoveryInvalid struct {
		c configuration.Provider
		m *RecoveryInvalidModel
	}
	RecoveryInvalidModel struct {
		To string
	}
)

func NewRecoveryInvalid(c configuration.Provider, m *RecoveryInvalidModel) *RecoveryInvalid {
	return &RecoveryInvalid{c: c, m: m}
}

func (t *RecoveryInvalid) EmailRecipient() (string, error) {
	return t.m.To, nil
}

func (t *RecoveryInvalid) EmailSubject() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "recovery/invalid/email.subject.gotmpl"), t.m)
}

func (t *RecoveryInvalid) EmailBody() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "recovery/invalid/email.body.gotmpl"), t.m)
}
