package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// ThawCmd implements `dynctl zone thaw`
var ThawCmd = &cobra.Command{
	Use:   "thaw <name>",
	Short: "Allow changes to again be made to a zone.",
	Long:  `Allow changes to again be made to a Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		if err := c.ThawZone(zone); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(zone, "thawed")
	},
}
