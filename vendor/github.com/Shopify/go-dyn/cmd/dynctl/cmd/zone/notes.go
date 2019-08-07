package zone

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/go-dyn/cmd/dynctl/cmd/util"
	"github.com/spf13/cobra"
)

// NotesCmd implements `dynctl zone notes`
var NotesCmd = &cobra.Command{
	Use:   "notes <name>",
	Short: "Retrieve notes for a single zone.",
	Long:  `Retrieve Zone Notes for a single Dyn Managed DNS zone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := util.DynLogIn()
		defer c.LogOut()

		zone := args[0]

		notes, err := c.GetZoneNotes(zone)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		fmt.Printf("%-18s  %-8s  %-26s  %s\n", "User", "Type", "When", "Note")

		for _, n := range notes {
			ts, _ := strconv.ParseInt(n.Timestamp, 10, 64)
			t := time.Unix(ts, 0).UTC()
			lines := strings.SplitAfter(n.Note, "\n")
			fmt.Printf("%-18s  %-8s  %-26s  ", n.UserName, n.Type, t.Format("Jan 02, 2006 (03:04 - MST)"))
			fmt.Print(lines[0])
			for _, l := range lines[1 : len(lines)-1] {
				fmt.Printf("%56s  %s", "", l)
			}
		}
	},
}
