// SPDX-License-Identifier: BSD-3-Clause
package set

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	setCmd := &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "Set or update object attributes.",
	}

	setCmd.AddCommand(userCmd)

	return setCmd
}
