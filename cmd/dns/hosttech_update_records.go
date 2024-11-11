package dns

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/DEXPRO-Solutions-GmbH/xd/pkg/dnsutils"
	"github.com/fatih/color"
	"github.com/libdns/hosttech"
	"github.com/manifoldco/promptui"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func newHosttechUpdateRecordsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-records",
	}

	flags := cmd.Flags()

	zone := flags.String("zone", "", "The zone to update")
	recordType := flags.String("record-type", "A", "The record type to update")
	domainStr := flags.String("domain", "", "Regular Expression: The record name to update")
	newValue := flags.String("new-value", "", "The new value to set")
	apiToken := flags.String("api-token", "", "The API token to use")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		domain, err := regexp.Compile(*domainStr)
		if err != nil {
			cmd.PrintErrf("invalid regular expression in domain: %v\n", err)
			os.Exit(1)
		}

		cmd.Printf("Using regular expression for domain matching: %v\n", domain.String())

		provider := &hosttech.Provider{
			APIToken: *apiToken,
		}

		records, err := dnsutils.PlanUpdate(context.Background(), domain, *recordType, *zone, *newValue, provider)
		if err != nil {
			cmd.PrintErrf("failed to plan update: %v\n", err)
			os.Exit(1)
		}

		cmd.Printf("This operation will update change the following %d records\n", len(records))

		if len(records) == 0 {
			cmd.Println("No records to update")
			os.Exit(0)
		}

		recordTable := table.New("Type", "Name", "Full Domain", "Current Value", "New Value")
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		recordTable.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, record := range records {
			recordTable.AddRow(record.Type, record.Name, fmt.Sprintf("%s.%s", record.Name, *zone), record.Value, *newValue)
		}

		cmd.Println()
		recordTable.Print()
		cmd.Println()

		// Prompt for user confirmation
		prompt := promptui.Prompt{
			Label:     color.RedString("Do you want to continue"),
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			os.Exit(1)
		}

		// User has accepted
		err = dnsutils.Update(context.Background(), *zone, records, *newValue, provider)
		if err != nil {
			cmd.PrintErrf("failed to update records: %v\n", err)
			os.Exit(1)
		}
	}

	return cmd
}
