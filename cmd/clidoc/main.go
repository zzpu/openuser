package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zzpu/openuser/cmd/remote"

	"github.com/ory/x/clidoc"

	"github.com/zzpu/openuser/cmd/identities"
	"github.com/zzpu/openuser/cmd/jsonnet"
)

func main() {
	rootCmd := &cobra.Command{Use: "kratos"}
	identities.RegisterCommandRecursive(rootCmd)
	jsonnet.RegisterCommandRecursive(rootCmd)
	remote.RegisterCommandRecursive(rootCmd)

	if err := clidoc.Generate(rootCmd, os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
	fmt.Println("All files have been generated and updated.")
}
