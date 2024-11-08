// SPDX-License-Identifier: BSD-3-Clause
package reset

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	resetCmd := &cobra.Command{
		Use:     "reset",
		Aliases: []string{"r"},
		Short:   "Reset objects.",
	}

	resetCmd.AddCommand(chassisCmd)
	resetCmd.AddCommand(systemCmd)

	return resetCmd
}
