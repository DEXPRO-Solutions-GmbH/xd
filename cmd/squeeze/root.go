package squeeze

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "squeeze",
	}

	cmd.AddCommand(newBenchmarkCmd())

	return cmd
}
