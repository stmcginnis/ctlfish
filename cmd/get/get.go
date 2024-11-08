// SPDX-License-Identifier: BSD-3-Clause
package get

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get object information.",
	}

	getCmd.AddCommand(chassisCmd)
	getCmd.AddCommand(driveCmd)
	getCmd.AddCommand(systemCmd)
	getCmd.AddCommand(userCmd)

	return getCmd
}
