package cmd

import (
	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/record"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(record.ListCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage records.",
	Long:  `Command to manage Dyn Managed DNS records.`,
}
