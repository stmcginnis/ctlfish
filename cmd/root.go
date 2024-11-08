// SPDX-License-Identifier: BSD-3-Clause
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/cmd/get"
	"github.com/stmcginnis/ctlfish/cmd/reset"
	"github.com/stmcginnis/ctlfish/cmd/set"
	"github.com/stmcginnis/ctlfish/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ctlfish",
	Short: "A Redfish and Swordfish CLI",
	Long:  `ctlfish is a CLI for interacting with Redfish and Swordfish systems.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctlfish.yaml)")
	config.InitConfig(cfgFile)

	rootCmd.AddCommand(get.Cmd())
	rootCmd.AddCommand(reset.Cmd())
	rootCmd.AddCommand(set.Cmd())
}
