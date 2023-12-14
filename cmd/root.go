package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "xd",
		Short: "DEXPRO Developer CLI",
	}

	return cmd
}
