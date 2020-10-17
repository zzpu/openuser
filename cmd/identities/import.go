package identities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zzpu/openuser/internal/clihelpers"

	"github.com/zzpu/openuser/cmd/cliclient"
	"github.com/zzpu/openuser/internal/httpclient/client/admin"
	"github.com/zzpu/openuser/internal/httpclient/models"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <file.json [file-2.json [file-3.json] ...]>",
	Short: "Import identities from files or STD_IN",
	Long: `Import identities from files or STD_IN:

	kratos identities import file.json

	cat file.json | kratos identities import

Files can contain only a single or an array of identities. The validity of files can be tested beforehand using "... identities validate".

WARNING: Importing credentials is not yet supported.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := cliclient.NewClient(cmd)

		imported := make([]*models.Identity, 0, len(args))
		failed := make(map[string]error)

		is, err := readIdentities(cmd, args)
		if err != nil {
			return err
		}

		for src, i := range is {
			err = validateIdentity(cmd, src, i, c.Public.GetSchema)
			if err != nil {
				return err
			}

			var params models.CreateIdentity
			err = json.Unmarshal([]byte(i), &params)
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "STD_IN: Could not parse identity")
				return clihelpers.FailSilently(cmd)
			}

			resp, err := c.Admin.CreateIdentity(&admin.CreateIdentityParams{
				Body:    &params,
				Context: context.Background(),
			})
			if err != nil {
				failed[src] = err
			} else {
				imported = append(imported, resp.Payload)
			}
		}
		if len(imported) == 1 {
			clihelpers.PrintRow(cmd, (*outputIdentity)(imported[0]))
		} else {
			clihelpers.PrintCollection(cmd, &outputIdentityCollection{identities: imported})
		}
		clihelpers.PrintErrors(cmd, failed)

		if len(failed) != 0 {
			return clihelpers.FailSilently(cmd)
		}

		return nil
	},
}
