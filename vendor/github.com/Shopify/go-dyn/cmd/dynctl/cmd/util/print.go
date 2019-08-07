package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// PrintJSON pretty prints the JSON representation of an object
func PrintJSON(i interface{}) error {
	e := json.NewEncoder(os.Stdout)

	e.SetIndent("", "\t")

	if err := e.Encode(i); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}
