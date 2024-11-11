package dns

import "github.com/spf13/cobra"

func newHosttechCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "hosttech",
	}

	cmd.AddCommand(newHosttechUpdateRecordsCommand())

	return cmd
}
