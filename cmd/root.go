package cmd

import (
	"github.com/DEXPRO-Solutions-GmbH/xd/cmd/squeeze"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xd",
		Short: "DEXPRO Developer CLI with an X",
	}

	cmd.AddCommand(newGenCmd())
	cmd.AddCommand(squeeze.NewRootCmd())

	return cmd
}
