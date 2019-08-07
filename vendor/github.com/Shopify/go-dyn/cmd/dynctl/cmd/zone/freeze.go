package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// FreezeCmd implements `dynctl zone freeze`
var FreezeCmd = &cobra.Command{
	Use:   "freeze <name>",
	Short: "Prevent changes to a zone until it is thawed.",
	Long:  `Prevent changes to a Dyn Managed DNS zone until it is thawed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		err := c.FreezeZone(zone)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(zone, "frozen")
	},
}
