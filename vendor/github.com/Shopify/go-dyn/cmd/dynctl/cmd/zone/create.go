package zone

import (
	"fmt"
	"os"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/cobra"
)

var rname string
var ttl int
var serialStyle serialStyleValue

var validSerialStyles = []string{
	dyn.SerialStyleIncrement,
	dyn.SerialStyleEpoch,
	dyn.SerialStyleDay,
	dyn.SerialStyleMinute,
}

type serialStyleValue struct {
	value string
}

func serialStyleValueValid(v string) bool {
	for _, s := range validSerialStyles {
		if v == s {
			return true
		}
	}

	return false
}

func (s *serialStyleValue) String() string {
	return s.value
}

func (s *serialStyleValue) Set(v string) error {
	if !serialStyleValueValid(v) {
		return fmt.Errorf("must be one of: %q", validSerialStyles)
	}

	s.value = v

	return nil
}

func (s *serialStyleValue) Type() string {
	return "serialStyle"
}

func init() {
	CreateCmd.Flags().StringVarP(&rname, "rname", "r", "", "administrative contact for this zone (required)")
	CreateCmd.MarkFlagRequired("rname")

	CreateCmd.Flags().VarP(&serialStyle, "serialstyle", "s", "style of the zone’s serial")

	CreateCmd.Flags().IntVarP(&ttl, "ttl", "t", 0, "default TTL (in seconds) for records in the zone (required)")
	CreateCmd.MarkFlagRequired("ttl")
}

// CreateCmd implements `dynctl zone create`
var CreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a primary zone.",
	Long:  `Create a Dyn Managed DNS primary zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		z, err := c.CreateZone(zone, rname, ttl, dyn.SerialStyle(serialStyle.value))

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(-1)
		}

		fmt.Println(zone, "created")
		util.PrintJSON(z)
	},
}
