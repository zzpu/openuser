package template

import (
	"path/filepath"

	"github.com/zzpu/openuser/driver/configuration"
)

type TestStub struct {
	c configuration.Provider
	m *TestStubModel
}

type TestStubModel struct {
	To      string
	Subject string
	Body    string
}

func NewTestStub(c configuration.Provider, m *TestStubModel) *TestStub {
	return &TestStub{c: c, m: m}
}

func (t *TestStub) EmailRecipient() (string, error) {
	return t.m.To, nil
}

func (t *TestStub) EmailSubject() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "test_stub/email.subject.gotmpl"), t.m)
}

func (t *TestStub) EmailBody() (string, error) {
	return loadTextTemplate(filepath.Join(t.c.CourierTemplatesRoot(), "test_stub/email.body.gotmpl"), t.m)
}
