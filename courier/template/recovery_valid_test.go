package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zzpu/openuser/courier/template"
	"github.com/zzpu/openuser/internal"
)

func TestRecoverValid(t *testing.T) {
	conf, _ := internal.NewFastRegistryWithMocks(t)
	tpl := template.NewRecoveryValid(conf, &template.RecoveryValidModel{})

	rendered, err := tpl.EmailBody()
	require.NoError(t, err)
	assert.NotEmpty(t, rendered)

	rendered, err = tpl.EmailSubject()
	require.NoError(t, err)
	assert.NotEmpty(t, rendered)
}
