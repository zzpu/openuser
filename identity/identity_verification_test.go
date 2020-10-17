package identity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zzpu/openuser/x"
)

func TestNewVerifiableEmailAddress(t *testing.T) {
	iid := x.NewUUID()
	a := NewVerifiableEmailAddress("foo@ory.sh", iid)

	assert.Equal(t, a.Value, "foo@ory.sh")
	assert.Equal(t, a.Via, VerifiableAddressTypeEmail)
	assert.Equal(t, a.Status, VerifiableAddressStatusPending)
	assert.Equal(t, a.Verified, false)
	assert.EqualValues(t, time.Time{}, a.VerifiedAt)
	assert.NotEmpty(t, a.ID)
}
