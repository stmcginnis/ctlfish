// SPDX-License-Identifier: BSD-3-Clause
package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/utils"
)

var systemCmd = &cobra.Command{
	Use:     "system [NAME_OR_ID]",
	Aliases: []string{"systems", "s"},
	Short:   "Get system information.",
	RunE:    getSystem,
	Args:    cobra.MaximumNArgs(1),
}

// getSystem retrieves the system information.
func getSystem(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	systems, err := c.Service.Systems()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve system information: %v", err)
	}

	writer := utils.NewTableWriter(
		cmd.OutOrStdout(),
		"name", "cpu", "memory", "power", "status", "led", "description")
	for _, system := range systems {
		if len(args) == 1 && (system.ID != args[0] && system.Name != args[0]) {
			continue
		}

		writer.AddRow(
			system.Name,
			system.ProcessorSummary.Count,
			fmt.Sprintf("%0.2f GB", system.MemorySummary.TotalSystemMemoryGiB),
			system.PowerState,
			system.Status.Health,
			system.IndicatorLED,
			system.Description)
	}

	if len(args) != 0 && writer.RowCount() == 0 {
		return utils.ErrorExit(cmd, "system '%s' was not found.", args[0])
	}

	writer.Render()
	return nil
}
