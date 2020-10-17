package identities

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/zzpu/openuser/identity"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/ory/viper"
	"github.com/zzpu/openuser/cmd/cliclient"
	"github.com/zzpu/openuser/driver"
	"github.com/zzpu/openuser/driver/configuration"
	"github.com/zzpu/openuser/internal"
	"github.com/zzpu/openuser/internal/clihelpers"
	"github.com/zzpu/openuser/internal/testhelpers"
)

func parseIdentities(raw []byte) (rawIdentities []string) {
	res := gjson.ParseBytes(raw)
	if !res.IsArray() {
		return []string{res.Raw}
	}
	res.ForEach(func(_, v gjson.Result) bool {
		rawIdentities = append(rawIdentities, v.Raw)
		return true
	})
	return
}

func readIdentities(cmd *cobra.Command, args []string) (map[string]string, error) {
	rawIdentities := make(map[string]string)
	if len(args) == 0 {
		fc, err := ioutil.ReadAll(cmd.InOrStdin())
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "STD_IN: Could not read: %s\n", err)
			return nil, clihelpers.FailSilently(cmd)
		}
		for i, id := range parseIdentities(fc) {
			rawIdentities[fmt.Sprintf("STD_IN[%d]", i)] = id
		}
		return rawIdentities, nil
	}
	for _, fn := range args {
		fc, err := ioutil.ReadFile(fn)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: Could not open identity file: %s\n", fn, err)
			return nil, clihelpers.FailSilently(cmd)
		}
		for i, id := range parseIdentities(fc) {
			rawIdentities[fmt.Sprintf("%s[%d]", fn, i)] = id
		}
	}
	return rawIdentities, nil
}

func setup(t *testing.T, cmd *cobra.Command) driver.Registry {
	_, reg := internal.NewRegistryDefaultWithDSN(t, configuration.DefaultSQLiteMemoryDSN)
	_, admin := testhelpers.NewKratosServerWithCSRF(t, reg)
	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stubs/identity.schema.json")
	// setup command
	cliclient.RegisterClientFlags(cmd.Flags())
	clihelpers.RegisterFormatFlags(cmd.Flags())
	require.NoError(t, cmd.Flags().Set(cliclient.FlagEndpoint, admin.URL))
	require.NoError(t, cmd.Flags().Set(clihelpers.FlagFormat, string(clihelpers.FormatJSON)))
	return reg
}

func exec(cmd *cobra.Command, stdIn io.Reader, args ...string) (string, string, error) {
	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.SetErr(stdErr)
	cmd.SetOut(stdOut)
	cmd.SetIn(stdIn)
	defer cmd.SetIn(nil)
	if args == nil {
		args = []string{}
	}
	cmd.SetArgs(args)
	err := cmd.Execute()
	return stdOut.String(), stdErr.String(), err
}

func execNoErr(t *testing.T, cmd *cobra.Command, args ...string) string {
	stdOut, stdErr, err := exec(cmd, nil, args...)
	require.NoError(t, err)
	require.Len(t, stdErr, 0, stdOut)
	return stdOut
}

func execErr(t *testing.T, cmd *cobra.Command, args ...string) string {
	stdOut, stdErr, err := exec(cmd, nil, args...)
	require.True(t, errors.Is(err, clihelpers.NoPrintButFailError))
	require.Len(t, stdOut, 0, stdErr)
	return stdErr
}

func makeIdentities(t *testing.T, reg driver.Registry, n int) (is []*identity.Identity, ids []string) {
	for j := 0; j < n; j++ {
		i := identity.NewIdentity(configuration.DefaultIdentityTraitsSchemaID)
		require.NoError(t, reg.Persister().CreateIdentity(context.Background(), i))
		is = append(is, i)
		ids = append(ids, i.ID.String())
	}
	return
}
