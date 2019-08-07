package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/cobra"
)

// ListCmd implements `dynctl zone list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available zones.",
	Long:  `List all available Dyn Managed DNS zones.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		n, err := c.EachZone(func(z *dyn.Zone) {
			util.PrintJSON(z)
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		fmt.Println(n, "zone(s)")
	},
}
