package link

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/stringslice"
	"github.com/ory/x/urlx"

	"github.com/zzpu/ums/selfservice/flow"
	"github.com/zzpu/ums/selfservice/flow/verification"
)

func TestVerificationToken(t *testing.T) {
	req := &http.Request{URL: urlx.ParseOrPanic("https://www.ory.sh/")}
	t.Run("func=NewSelfServiceVerificationToken", func(t *testing.T) {
		t.Run("case=creates unique tokens", func(t *testing.T) {
			f, err := verification.NewFlow(time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			tokens := make([]string, 10)
			for k := range tokens {
				tokens[k] = NewSelfServiceVerificationToken(nil, f).Token
			}

			assert.Len(t, stringslice.Unique(tokens), len(tokens))
		})
	})
	t.Run("method=Valid", func(t *testing.T) {
		t.Run("case=is invalid when the flow is expired", func(t *testing.T) {
			f, err := verification.NewFlow(-time.Hour, "", req, nil, flow.TypeBrowser)
			require.NoError(t, err)

			token := NewSelfServiceVerificationToken(nil, f)
			require.Error(t, token.Valid())
			assert.EqualError(t, token.Valid(), f.Valid().Error())
		})
	})
}
