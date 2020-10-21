package identities

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlcon"
	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/x"
)

func TestDeleteCmd(t *testing.T) {
	reg := setup(t, deleteCmd)

	t.Run("case=deletes successfully", func(t *testing.T) {
		// create identity to delete
		i := identity.NewIdentity(configuration.DefaultIdentityTraitsSchemaID)
		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))

		stdOut := execNoErr(t, deleteCmd, i.ID.String())

		// expect ID and no error
		assert.Equal(t, i.ID.String()+"\n", stdOut)

		// expect identity to be deleted
		_, err := reg.Persister().GetIdentity(context.Background(), i.ID)
		assert.True(t, errors.Is(err, sqlcon.ErrNoRows))
	})

	t.Run("case=deletes three identities", func(t *testing.T) {
		is, ids := makeIdentities(t, reg, 3)

		stdOut := execNoErr(t, deleteCmd, ids...)

		assert.Equal(t, strings.Join(ids, "\n")+"\n", stdOut)

		for _, i := range is {
			_, err := reg.Persister().GetIdentity(context.Background(), i.ID)
			assert.Error(t, err)
		}
	})

	t.Run("case=fails with unknown ID", func(t *testing.T) {
		stdErr := execErr(t, deleteCmd, x.NewUUID().String())

		assert.Contains(t, stdErr, "[DELETE /identities/{id}][404] deleteIdentityNotFound", stdErr)
	})
}
