package td

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/cobra"
)

var label string

func init() {
	ListCmd.Flags().StringVarP(&label, "label", "l", "", "list service(s) with specified label (can be a wildcard)")
}

// ListCmd implements `dynctl td list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Traffic Director service instances.",
	Long:  `List all Dyn Managed DNS Traffic Director service instances.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		n, err := c.EachTrafficDirector(func(td *dyn.TrafficDirector) {
			util.PrintJSON(td)
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		fmt.Println(n, "Traffic Director service instance(s)")
	},
}
