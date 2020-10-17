package link

import (
	"context"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"
	"github.com/ory/x/assertx"

	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/identity"
	"github.com/zzpu/openuser/selfservice/flow/recovery"
	"github.com/zzpu/openuser/selfservice/flow/verification"
	"github.com/zzpu/openuser/x"
)

func TestPersister(p interface {
	RecoveryTokenPersister
	VerificationTokenPersister
	recovery.FlowPersister
	verification.FlowPersister
	identity.PrivilegedPool
}) func(t *testing.T) {
	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stub/identity.schema.json")
	viper.Set(configuration.ViperKeySecretsDefault, []string{"secret-a", "secret-b"})
	return func(t *testing.T) {
		t.Run("token=recovery", func(t *testing.T) {

			t.Run("case=should error when the recovery token does not exist", func(t *testing.T) {
				_, err := p.UseRecoveryToken(context.Background(), "i-do-not-exist")
				require.Error(t, err)
			})

			newRecoveryToken := func(t *testing.T, email string) *RecoveryToken {
				var req recovery.Flow
				require.NoError(t, faker.FakeData(&req))
				require.NoError(t, p.CreateRecoveryFlow(context.Background(), &req))

				var i identity.Identity
				require.NoError(t, faker.FakeData(&i))

				address := &identity.RecoveryAddress{Value: email, Via: identity.RecoveryAddressTypeEmail}
				i.RecoveryAddresses = append(i.RecoveryAddresses, *address)

				require.NoError(t, p.CreateIdentity(context.Background(), &i))

				return &RecoveryToken{Token: x.NewUUID().String(), FlowID: uuid.NullUUID{UUID: req.ID, Valid: true},
					RecoveryAddress: &i.RecoveryAddresses[0],
					ExpiresAt:       time.Now(),
					IssuedAt:        time.Now(),
				}
			}

			t.Run("case=should error when the recovery token does not exist", func(t *testing.T) {
				_, err := p.UseRecoveryToken(context.Background(), "i-do-not-exist")
				require.Error(t, err)
			})

			t.Run("case=should create a new recovery token", func(t *testing.T) {
				token := newRecoveryToken(t, "foo-user@ory.sh")
				require.NoError(t, p.CreateRecoveryToken(context.Background(), token))
			})

			t.Run("case=should create a recovery token and use it", func(t *testing.T) {
				expected := newRecoveryToken(t, "other-user@ory.sh")
				require.NoError(t, p.CreateRecoveryToken(context.Background(), expected))
				actual, err := p.UseRecoveryToken(context.Background(), expected.Token)
				require.NoError(t, err)
				assertx.EqualAsJSON(t, expected.RecoveryAddress, actual.RecoveryAddress)
				assertx.EqualAsJSON(t, expected.RecoveryAddress, actual.RecoveryAddress)
				assert.Equal(t, expected.RecoveryAddress.IdentityID, actual.RecoveryAddress.IdentityID)
				assert.NotEqual(t, expected.Token, actual.Token)
				assert.EqualValues(t, expected.FlowID, actual.FlowID)

				_, err = p.UseRecoveryToken(context.Background(), expected.Token)
				require.Error(t, err)
			})

		})
		t.Run("token=verification", func(t *testing.T) {

			t.Run("case=should error when the verification token does not exist", func(t *testing.T) {
				_, err := p.UseVerificationToken(context.Background(), "i-do-not-exist")
				require.Error(t, err)
			})

			newVerificationToken := func(t *testing.T, email string) *VerificationToken {
				var req verification.Flow
				require.NoError(t, faker.FakeData(&req))
				require.NoError(t, p.CreateVerificationFlow(context.Background(), &req))

				var i identity.Identity
				require.NoError(t, faker.FakeData(&i))

				address := &identity.VerifiableAddress{Value: email, Via: identity.VerifiableAddressTypeEmail}
				i.VerifiableAddresses = append(i.VerifiableAddresses, *address)

				require.NoError(t, p.CreateIdentity(context.Background(), &i))

				return &VerificationToken{
					Token:             x.NewUUID().String(),
					FlowID:            uuid.NullUUID{UUID: req.ID, Valid: true},
					VerifiableAddress: &i.VerifiableAddresses[0],
					ExpiresAt:         time.Now(),
					IssuedAt:          time.Now(),
				}
			}

			t.Run("case=should error when the verification token does not exist", func(t *testing.T) {
				_, err := p.UseVerificationToken(context.Background(), "i-do-not-exist")
				require.Error(t, err)
			})

			t.Run("case=should create a new verification token", func(t *testing.T) {
				token := newVerificationToken(t, "foo-user@ory.sh")
				require.NoError(t, p.CreateVerificationToken(context.Background(), token))
			})

			t.Run("case=should create a verification token and use it", func(t *testing.T) {
				expected := newVerificationToken(t, "other-user@ory.sh")
				require.NoError(t, p.CreateVerificationToken(context.Background(), expected))
				actual, err := p.UseVerificationToken(context.Background(), expected.Token)
				require.NoError(t, err)
				assertx.EqualAsJSON(t, expected.VerifiableAddress, actual.VerifiableAddress)
				assertx.EqualAsJSON(t, expected.VerifiableAddress, actual.VerifiableAddress)
				assert.Equal(t, expected.VerifiableAddress.IdentityID, actual.VerifiableAddress.IdentityID)
				assert.NotEqual(t, expected.Token, actual.Token)
				assert.EqualValues(t, expected.FlowID, actual.FlowID)

				_, err = p.UseVerificationToken(context.Background(), expected.Token)
				require.Error(t, err)
			})
		})
	}
}
