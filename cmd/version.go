//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"github.com/spf13/cobra"
)

// BuildVersion is set in official builds.
var BuildVersion string = "v0.0.0-dev"

// versionCmd prints the version of the CLI.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "ctlfish build version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(BuildVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
