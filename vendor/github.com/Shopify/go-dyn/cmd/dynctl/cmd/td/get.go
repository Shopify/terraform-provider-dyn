package td

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// GetCmd implements `dynctl td get`
var GetCmd = &cobra.Command{
	Use:   "get <label>",
	Short: "Retrieve information about a Traffic Director instance.",
	Long:  `Retrieve information about a Dyn Managed DNS Traffic Director instance.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		label := args[0]

		td, err := c.FindTrafficDirector(label)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		if err = util.PrintJSON(td); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	},
}
