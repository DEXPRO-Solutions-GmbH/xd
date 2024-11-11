package dns

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dns",
	}

	cmd.AddCommand(newHosttechCommand())

	return cmd
}
