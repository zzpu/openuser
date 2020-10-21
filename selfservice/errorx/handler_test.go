package errorx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/nosurf"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/errorsx"

	"github.com/ory/herodot"

	"github.com/zzpu/ums/internal"
	"github.com/zzpu/ums/selfservice/errorx"
	"github.com/zzpu/ums/x"
)

func TestHandler(t *testing.T) {
	_, reg := internal.NewFastRegistryWithMocks(t)
	h := errorx.NewHandler(reg)

	t.Run("case=public authorization", func(t *testing.T) {
		router := x.NewRouterPublic()
		ns := x.NewTestCSRFHandler(router, reg)

		h.RegisterPublicRoutes(router)
		router.GET("/regen", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			ns.RegenerateToken(w, r)
			w.WriteHeader(http.StatusNoContent)
		})
		router.GET("/set-error", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			id, err := reg.SelfServiceErrorPersister().Add(context.Background(), nosurf.Token(r), herodot.ErrNotFound.WithReason("foobar"))
			require.NoError(t, err)
			_, _ = w.Write([]byte(id.String()))
		})

		ts := httptest.NewServer(ns)
		defer ts.Close()

		getBody := func(t *testing.T, hc *http.Client, path string, expectedCode int) []byte {
			res, err := hc.Get(ts.URL + path)
			require.NoError(t, err)
			defer res.Body.Close()
			require.EqualValues(t, expectedCode, res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			return body
		}
		expectedError := x.MustEncodeJSON(t, []error{herodot.ErrNotFound.WithReason("foobar")})

		t.Run("call with valid csrf cookie", func(t *testing.T) {
			hc := &http.Client{}
			id := getBody(t, hc, "/set-error", http.StatusOK)
			actual := getBody(t, hc, errorx.RouteGet+"?error="+string(id), http.StatusOK)
			assert.JSONEq(t, expectedError, gjson.GetBytes(actual, "errors").Raw, "%s", actual)

			// We expect a forbid error if the error is not found, regardless of CSRF
			_ = getBody(t, hc, errorx.RouteGet+"?error=does-not-exist", http.StatusForbidden)
		})
	})

	t.Run("case=stubs", func(t *testing.T) {
		router := x.NewRouterAdmin()
		h.RegisterAdminRoutes(router)
		ts := httptest.NewServer(router)
		defer ts.Close()

		res, err := ts.Client().Get(ts.URL + errorx.RouteGet + "?error=stub:500")
		require.NoError(t, err)
		require.EqualValues(t, http.StatusOK, res.StatusCode)

		actual, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.EqualValues(t, "This is a stub error.", gjson.GetBytes(actual, "errors.0.reason").String())
	})

	t.Run("case=errors types", func(t *testing.T) {
		router := x.NewRouterAdmin()
		h.RegisterAdminRoutes(router)
		ts := httptest.NewServer(router)
		defer ts.Close()

		for k, tc := range []struct {
			gave []error
		}{
			{
				gave: []error{
					herodot.ErrNotFound.WithReason("foobar"),
				},
			},
			{
				gave: []error{
					herodot.ErrNotFound.WithReason("foobar"),
					herodot.ErrNotFound.WithReason("foobar"),
				},
			},
			{
				gave: []error{
					herodot.ErrNotFound.WithReason("foobar"),
				},
			},
			{
				gave: []error{
					errors.WithStack(herodot.ErrNotFound.WithReason("foobar")),
				},
			},
			{
				gave: []error{
					errors.WithStack(herodot.ErrNotFound.WithReason("foobar").WithTrace(errors.New("asdf"))),
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				csrf := x.NewUUID()
				id, err := reg.SelfServiceErrorPersister().Add(context.Background(), csrf.String(), tc.gave...)
				require.NoError(t, err)

				res, err := ts.Client().Get(ts.URL + errorx.RouteGet + "?error=" + id.String())
				require.NoError(t, err)
				defer res.Body.Close()
				assert.EqualValues(t, http.StatusOK, res.StatusCode)

				actual, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)

				gg := make([]error, len(tc.gave))
				for k, g := range tc.gave {
					gg[k] = errorsx.Cause(g)
				}

				expected, err := json.Marshal(errorx.ErrorContainer{
					ID:     id,
					Errors: x.RequireJSONMarshal(t, gg),
				})
				require.NoError(t, err)

				assert.JSONEq(t, string(expected), string(actual), "%s != %s", expected, actual)
				assert.Empty(t, gjson.GetBytes(actual, "csrf_token").String())
				assert.JSONEq(t, string(x.RequireJSONMarshal(t, gg)), gjson.GetBytes(actual, "errors").Raw)
				t.Logf("%s", actual)
			})
		}
	})
}
