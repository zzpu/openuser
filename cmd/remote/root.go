package remote

import (
	"github.com/spf13/cobra"

	"github.com/zzpu/openuser/cmd/cliclient"
	"github.com/zzpu/openuser/internal/clihelpers"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Helpers and management for remote ORY Kratos instances",
}

func RegisterCommandRecursive(parent *cobra.Command) {
	parent.AddCommand(remoteCmd)

	remoteCmd.AddCommand(versionCmd)
	remoteCmd.AddCommand(statusCmd)
}

func init() {
	cliclient.RegisterClientFlags(remoteCmd.PersistentFlags())
	clihelpers.RegisterFormatFlags(remoteCmd.PersistentFlags())
}
