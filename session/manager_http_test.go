package session_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/internal"
	"github.com/zzpu/ums/internal/testhelpers"
	"github.com/zzpu/ums/session"
	"github.com/zzpu/ums/x"
)

var _ x.CSRFHandler = new(mockCSRFHandler)

type mockCSRFHandler struct {
	c int
}

func (f *mockCSRFHandler) ExemptPath(s string) {
}

func (f *mockCSRFHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (f *mockCSRFHandler) RegenerateToken(w http.ResponseWriter, r *http.Request) string {
	f.c++
	return x.FakeCSRFToken
}

func TestManagerHTTP(t *testing.T) {
	t.Run("case=regenerate csrf on principal change", func(t *testing.T) {
		_, reg := internal.NewFastRegistryWithMocks(t)
		mock := new(mockCSRFHandler)
		reg.WithCSRFHandler(mock)

		require.NoError(t, reg.SessionManager().IssueCookie(context.Background(), httptest.NewRecorder(), new(http.Request), new(session.Session)))
		assert.Equal(t, 1, mock.c)
	})

	t.Run("suite=lifecycle", func(t *testing.T) {
		conf, reg := internal.NewFastRegistryWithMocks(t)

		viper.Set(configuration.ViperKeySelfServiceLoginUI, "https://www.ory.sh")
		viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stub/fake-session.schema.json")

		var s *session.Session
		rp := x.NewRouterPublic()
		rp.GET("/session/revoke", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			require.NoError(t, reg.SessionManager().PurgeFromRequest(r.Context(), w, r))
			w.WriteHeader(http.StatusOK)
		})

		rp.GET("/session/set", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			require.NoError(t, reg.SessionManager().CreateAndIssueCookie(r.Context(), w, r, s))
			w.WriteHeader(http.StatusOK)
		})

		rp.GET("/session/get", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			sess, err := reg.SessionManager().FetchFromRequest(r.Context(), r)
			if err != nil {
				t.Logf("Got error on lookup: %s %T", err, errors.Unwrap(err))
				reg.Writer().WriteError(w, r, err)
				return
			}
			reg.Writer().Write(w, r, sess)
		})

		pts := httptest.NewServer(x.NewTestCSRFHandler(rp, reg))
		t.Cleanup(pts.Close)
		viper.Set(configuration.ViperKeyPublicBaseURL, pts.URL)
		reg.RegisterPublicRoutes(rp)

		t.Run("case=valid", func(t *testing.T) {
			viper.Set(configuration.ViperKeySessionLifespan, "1m")

			i := identity.Identity{Traits: []byte("{}")}
			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), &i))
			s = session.NewActiveSession(&i, conf, time.Now())

			c := testhelpers.NewClientWithCookies(t)
			testhelpers.MockHydrateCookieClient(t, c, pts.URL+"/session/set")

			res, err := c.Get(pts.URL + "/session/get")
			require.NoError(t, err)
			assert.EqualValues(t, http.StatusOK, res.StatusCode)
		})

		t.Run("case=expired", func(t *testing.T) {
			viper.Set(configuration.ViperKeySessionLifespan, "1ns")
			defer viper.Set(configuration.ViperKeySessionLifespan, "1m")

			i := identity.Identity{Traits: []byte("{}")}
			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), &i))
			s = session.NewActiveSession(&i, conf, time.Now())

			c := testhelpers.NewClientWithCookies(t)
			testhelpers.MockHydrateCookieClient(t, c, pts.URL+"/session/set")

			time.Sleep(time.Nanosecond * 2)

			res, err := c.Get(pts.URL + "/session/get")
			require.NoError(t, err)
			assert.EqualValues(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("case=revoked", func(t *testing.T) {
			i := identity.Identity{Traits: []byte("{}")}
			require.NoError(t, reg.PrivilegedIdentityPool().CreateIdentity(context.Background(), &i))
			s = session.NewActiveSession(&i, conf, time.Now())

			s = session.NewActiveSession(&i, conf, time.Now())

			c := testhelpers.NewClientWithCookies(t)
			testhelpers.MockHydrateCookieClient(t, c, pts.URL+"/session/set")

			res, err := c.Get(pts.URL + "/session/revoke")
			require.NoError(t, err)
			assert.EqualValues(t, http.StatusOK, res.StatusCode)

			res, err = c.Get(pts.URL + "/session/get")
			require.NoError(t, err)
			assert.EqualValues(t, http.StatusUnauthorized, res.StatusCode)
		})
	})
}
