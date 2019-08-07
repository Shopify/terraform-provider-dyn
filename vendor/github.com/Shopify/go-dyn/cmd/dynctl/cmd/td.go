package cmd

import (
	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/td"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tdCmd)
	tdCmd.AddCommand(td.CreateCmd)
	tdCmd.AddCommand(td.ListCmd)
	tdCmd.AddCommand(td.GetCmd)
}

var tdCmd = &cobra.Command{
	Use:   "td",
	Short: "Manage Traffic Director service instances.",
	Long:  `Command to manage Traffic Director service instances.`,
}
