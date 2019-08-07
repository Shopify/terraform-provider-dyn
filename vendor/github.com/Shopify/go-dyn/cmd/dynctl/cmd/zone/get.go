package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// GetCmd implements `dynctl zone get`
var GetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Retrieve information about a single zone.",
	Long:  `Retrieve information about a single Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		z, err := c.GetZone(zone)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		if err = util.PrintJSON(z); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	},
}
