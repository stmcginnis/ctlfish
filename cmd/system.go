//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
)

// systemCmd represents the system command
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for viewing and interacting with system objects",
}

func init() {
	rootCmd.AddCommand(systemCmd)
	systemCmd.PersistentFlags().StringP("connection", "c", config.GetDefault(), "The stored connection name to use.")
}
