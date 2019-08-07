package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

var notes string

func init() {
	PublishCmd.Flags().StringVarP(&notes, "notes", "n", "", "custom note to be added to zone notes")
}

// PublishCmd implements `dynctl zone publish`
var PublishCmd = &cobra.Command{
	Use:   "publish <name>",
	Short: "Cause all pending changes to become part of a zone.",
	Long:  `Cause all pending changes to become part of a Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		z, err := c.PublishZone(zone, notes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(zone, "published")
		util.PrintJSON(z)
	},
}
