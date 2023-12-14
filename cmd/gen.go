package cmd

import "github.com/spf13/cobra"

func newGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate stuff",
	}

	cmd.AddCommand(newGenProjectIDCommand())

	return cmd
}
