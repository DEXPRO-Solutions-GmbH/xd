package cmd

import (
	"fmt"

	"github.com/DEXPRO-Solutions-GmbH/xd/pkg/dexpro"
	"github.com/spf13/cobra"
)

func newGenProjectIDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project-id <domain>",
		Short: "Generate a project ID based on a Domain",
		Long: `Generate a project ID based on a Domain. Result is printed to stdout.

The generated ID will always be the same since it is a UUIDv5 derived from a hardcoded namespace UUID.`,
		Args: cobra.ExactArgs(1),
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		domain := args[0]

		tenantID := dexpro.GenerateProjectID(domain)

		fmt.Println(tenantID)
	}

	return cmd
}
