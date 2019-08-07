package td

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/cobra"
)

var ttl int

func init() {
	CreateCmd.Flags().IntVarP(&ttl, "ttl", "t", 0, "default TTL (in seconds) to be used across the service (required)")
	CreateCmd.MarkFlagRequired("ttl")
}

// CreateCmd implements `dynctl td create`
var CreateCmd = &cobra.Command{
	Use:   "create <label>",
	Short: "Create a Traffic Director service instance.",
	Long:  `Create a Dyn Managed DNS Traffic Director service instance.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		label := args[0]
		optionsSetter := func(req *dyn.TrafficDirectorCURequest) {
			req.TTL = ttl
		}

		td, err := c.CreateTrafficDirector(label, optionsSetter)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(label, "created")
		util.PrintJSON(td)
	},
}
