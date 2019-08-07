package util

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/go-dyn/pkg/dyn"
	"github.com/spf13/viper"
)

// DynLogIn returns a logged in Dyn API client
func DynLogIn() *dyn.Client {
	c := dyn.NewClient()

	if viper.GetBool("verbose") {
		c.Logger = log.New(os.Stderr, "[dynctl] ", 0)
	}

	customer := viper.GetString("customer")
	user := viper.GetString("user")
	password := viper.GetString("password")

	if err := c.LogIn(customer, user, password); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	}

	return c
}
