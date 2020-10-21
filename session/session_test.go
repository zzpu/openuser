package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/internal"
	"github.com/zzpu/ums/session"
)

func TestSession(t *testing.T) {
	conf, _ := internal.NewFastRegistryWithMocks(t)
	authAt := time.Now()

	s := session.NewActiveSession(new(identity.Identity), conf, authAt)
	assert.True(t, s.IsActive())

	assert.False(t, (&session.Session{ExpiresAt: time.Now().Add(time.Hour)}).IsActive())
	assert.False(t, (&session.Session{Active: true}).IsActive())
}
