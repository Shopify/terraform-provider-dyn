package cmd

import (
	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/zone"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(zoneCmd)
	zoneCmd.AddCommand(zone.CreateCmd)
	zoneCmd.AddCommand(zone.ListCmd)
	zoneCmd.AddCommand(zone.GetCmd)
	zoneCmd.AddCommand(zone.DeleteCmd)
	zoneCmd.AddCommand(zone.PublishCmd)
	zoneCmd.AddCommand(zone.FreezeCmd)
	zoneCmd.AddCommand(zone.ThawCmd)
	zoneCmd.AddCommand(zone.NotesCmd)
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage zones.",
	Long:  `Command to manage Dyn Managed DNS zones.`,
}
