//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"github.com/spf13/cobra"
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
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctlfish.yaml)")
	config.InitConfig(cfgFile)
}
