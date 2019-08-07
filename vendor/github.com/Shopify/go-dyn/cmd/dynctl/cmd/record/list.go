package record

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/cobra"
)

// ListCmd implements `dynctl record list`
var ListCmd = &cobra.Command{
	Use:   "list <name>",
	Short: "List all records for the specified zone.",
	Long:  `List all records for the specified Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		n, err := c.EachRecord(zone, func(r *dyn.Record) { fmt.Println(r) })

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		fmt.Println(n, "record(s)")
	},
}
