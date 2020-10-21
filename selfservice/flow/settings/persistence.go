package settings

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/form"
	"github.com/zzpu/ums/x"
)

type (
	FlowPersister interface {
		CreateSettingsFlow(context.Context, *Flow) error
		GetSettingsFlow(ctx context.Context, id uuid.UUID) (*Flow, error)
		UpdateSettingsFlow(context.Context, *Flow) error
		UpdateSettingsFlowMethod(context.Context, uuid.UUID, string, *FlowMethod) error
	}
	FlowPersistenceProvider interface {
		SettingsFlowPersister() FlowPersister
	}
)

func TestRequestPersister(p interface {
	FlowPersister
	identity.PrivilegedPool
}) func(t *testing.T) {
	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stub/identity.schema.json")

	var clearids = func(r *Flow) {
		r.ID = uuid.UUID{}
		r.Identity.ID = uuid.UUID{}
		r.IdentityID = uuid.UUID{}
	}

	return func(t *testing.T) {
		t.Run("case=should error when the settings request does not exist", func(t *testing.T) {
			_, err := p.GetSettingsFlow(context.Background(), x.NewUUID())
			require.Error(t, err)
		})

		var newFlow = func(t *testing.T) *Flow {
			var r Flow
			require.NoError(t, faker.FakeData(&r))
			clearids(&r)
			require.NoError(t, p.CreateIdentity(context.Background(), r.Identity))
			return &r
		}

		t.Run("case=should create a new settings request", func(t *testing.T) {
			r := newFlow(t)
			err := p.CreateSettingsFlow(context.Background(), r)
			require.NoError(t, err, "%#v", err)
		})

		t.Run("case=should create with set ids", func(t *testing.T) {
			var r Flow
			require.NoError(t, faker.FakeData(&r))
			require.NoError(t, p.CreateIdentity(context.Background(), r.Identity))
			require.NoError(t, p.CreateSettingsFlow(context.Background(), &r))
		})

		t.Run("case=should create and fetch a settings request", func(t *testing.T) {
			expected := newFlow(t)
			err := p.CreateSettingsFlow(context.Background(), expected)
			require.NoError(t, err)

			actual, err := p.GetSettingsFlow(context.Background(), expected.ID)
			require.NoError(t, err)

			factual, _ := json.Marshal(actual.Methods[StrategyProfile].Config)
			fexpected, _ := json.Marshal(expected.Methods[StrategyProfile].Config)

			require.NotEmpty(t, actual.Methods[StrategyProfile].Config.FlowMethodConfigurator.(*form.HTMLForm).Action)
			assert.EqualValues(t, expected.ID, actual.ID)
			assert.JSONEq(t, string(fexpected), string(factual))
			x.AssertEqualTime(t, expected.IssuedAt, actual.IssuedAt)
			x.AssertEqualTime(t, expected.ExpiresAt, actual.ExpiresAt)
			assert.EqualValues(t, expected.RequestURL, actual.RequestURL)
			assert.EqualValues(t, expected.Identity.ID, actual.Identity.ID)
			assert.EqualValues(t, expected.Identity.Traits, actual.Identity.Traits)
			assert.EqualValues(t, expected.Identity.SchemaID, actual.Identity.SchemaID)
			assert.Empty(t, actual.Identity.Credentials)
		})

		t.Run("case=should fail to create if identity does not exist", func(t *testing.T) {
			var expected Flow
			require.NoError(t, faker.FakeData(&expected))
			clearids(&expected)
			expected.Identity = nil
			expected.IdentityID = uuid.Nil
			err := p.CreateSettingsFlow(context.Background(), &expected)
			require.Error(t, err, "%+s", expected)
		})

		t.Run("case=should create and update a settings request", func(t *testing.T) {
			expected := newFlow(t)
			expected.Methods["oidc"] = &FlowMethod{
				Method: "oidc", Config: &FlowMethodConfig{FlowMethodConfigurator: &form.HTMLForm{Fields: []form.Field{{
					Name: "zab", Type: "bar", Pattern: "baz"}}}}}
			expected.Methods["password"] = &FlowMethod{
				Method: "password", Config: &FlowMethodConfig{FlowMethodConfigurator: &form.HTMLForm{Fields: []form.Field{{
					Name: "foo", Type: "bar", Pattern: "baz"}}}}}
			err := p.CreateSettingsFlow(context.Background(), expected)
			require.NoError(t, err)

			expected.Methods[StrategyProfile].Config.FlowMethodConfigurator.(*form.HTMLForm).Action = "/new-action"
			expected.Methods["password"].Config.FlowMethodConfigurator.(*form.HTMLForm).Fields = []form.Field{{
				Name: "zab", Type: "zab", Pattern: "zab"}}
			expected.RequestURL = "/new-request-url"
			require.NoError(t, p.UpdateSettingsFlow(context.Background(), expected))

			actual, err := p.GetSettingsFlow(context.Background(), expected.ID)
			require.NoError(t, err)

			assert.Equal(t, "/new-action", actual.Methods[StrategyProfile].Config.FlowMethodConfigurator.(*form.HTMLForm).Action)
			assert.Equal(t, "/new-request-url", actual.RequestURL)
			assert.EqualValues(t, []form.Field{{Name: "zab", Type: "zab", Pattern: "zab"}}, actual.
				Methods["password"].Config.FlowMethodConfigurator.(*form.HTMLForm).Fields)
			assert.EqualValues(t, []form.Field{{Name: "zab", Type: "bar", Pattern: "baz"}}, actual.
				Methods["oidc"].Config.FlowMethodConfigurator.(*form.HTMLForm).Fields)
		})

		t.Run("case=should update a settings flow method", func(t *testing.T) {
			expected := newFlow(t)
			delete(expected.Methods, identity.CredentialsTypeOIDC.String())
			delete(expected.Methods, StrategyProfile)

			err := p.CreateSettingsFlow(context.Background(), expected)
			require.NoError(t, err)

			actual, err := p.GetSettingsFlow(context.Background(), expected.ID)
			require.NoError(t, err)
			assert.Len(t, actual.Methods, 1)

			require.NoError(t, p.UpdateSettingsFlowMethod(context.Background(), expected.ID, identity.CredentialsTypeOIDC.String(), &FlowMethod{
				Method: identity.CredentialsTypeOIDC.String(),
				Config: &FlowMethodConfig{FlowMethodConfigurator: form.NewHTMLForm(string(identity.CredentialsTypeOIDC))},
			}))

			require.NoError(t, p.UpdateSettingsFlowMethod(context.Background(), expected.ID, identity.CredentialsTypePassword.String(), &FlowMethod{
				Method: identity.CredentialsTypePassword.String(),
				Config: &FlowMethodConfig{FlowMethodConfigurator: form.NewHTMLForm(string(identity.CredentialsTypePassword))},
			}))

			actual, err = p.GetSettingsFlow(context.Background(), expected.ID)
			require.NoError(t, err)
			require.Len(t, actual.Methods, 2)
			assert.EqualValues(t, identity.CredentialsTypePassword, actual.Active)

			assert.Equal(t, string(identity.CredentialsTypePassword), actual.Methods[identity.CredentialsTypePassword.String()].Config.FlowMethodConfigurator.(*form.HTMLForm).Action)
			assert.Equal(t, string(identity.CredentialsTypeOIDC), actual.Methods[identity.CredentialsTypeOIDC.String()].Config.FlowMethodConfigurator.(*form.HTMLForm).Action)
		})
	}
}
