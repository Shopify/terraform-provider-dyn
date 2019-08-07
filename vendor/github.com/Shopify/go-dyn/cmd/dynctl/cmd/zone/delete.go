package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// DeleteCmd implements `dynctl zone delete`
var DeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Remove a single zone.",
	Long:  `Remove a Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		err := c.DeleteZone(zone)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(zone, "removed")
	},
}
