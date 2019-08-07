package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dynCustomerName, dynUserName, dynPassword string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "dynctl",
	Short: "Command-line utility to work with Dyn Managed DNS",
	Long:  `A Go package and CLI tool for working with Dyn Managed DNS.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dynctl.yml)")

	rootCmd.PersistentFlags().StringVar(&dynCustomerName, "customer", "", "Dyn customer name")
	viper.BindPFlag("customer", rootCmd.PersistentFlags().Lookup("customer"))
	viper.BindEnv("customer", "DYN_CUSTOMER_NAME")

	rootCmd.PersistentFlags().StringVar(&dynUserName, "user", "", "Dyn user name")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindEnv("user", "DYN_USER_NAME")

	rootCmd.PersistentFlags().StringVar(&dynPassword, "password", "", "Dyn password")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindEnv("password", "DYN_PASSWORD")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dynctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dynctl")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Can't read config:", err)
			os.Exit(1)
		}
	}
}

// Execute invokes the main dynctl command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
