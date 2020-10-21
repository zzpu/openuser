package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/zzpu/ums/cmd/remote"

	"github.com/ory/x/cmdx"
	"github.com/zzpu/ums/cmd/identities"
	"github.com/zzpu/ums/cmd/jsonnet"
	"github.com/zzpu/ums/cmd/migrate"
	"github.com/zzpu/ums/cmd/serve"
	"github.com/zzpu/ums/internal/clihelpers"

	"github.com/ory/x/viperx"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "kratos",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if !errors.Is(err, clihelpers.NoPrintButFailError) {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func init() {
	viperx.RegisterConfigFlag(rootCmd, "kratos")

	identities.RegisterCommandRecursive(rootCmd)
	jsonnet.RegisterCommandRecursive(rootCmd)
	serve.RegisterCommandRecursive(rootCmd)
	migrate.RegisterCommandRecursive(rootCmd)
	remote.RegisterCommandRecursive(rootCmd)

	rootCmd.AddCommand(cmdx.Version(&clihelpers.BuildVersion, &clihelpers.BuildGitHash, &clihelpers.BuildTime))
}
