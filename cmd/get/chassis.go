// SPDX-License-Identifier: BSD-3-Clause
package get

import (
	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/utils"
)

var chassisCmd = &cobra.Command{
	Use:     "chassis",
	Aliases: []string{"c"},
	Short:   "Get information about chassis objects.",
	RunE:    getChassis,
}

// getChassis retrieves the chassis information from the system.
func getChassis(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	chassis, err := c.Service.Chassis()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve chassis information: %v", err)
	}

	writer := utils.NewTableWriter(cmd.OutOrStdout(), "name", "power", "status")
	for _, chass := range chassis {
		if len(args) == 1 && (chass.ID != args[0] && chass.Name != args[0]) {
			continue
		}

		writer.AddRow(chass.Name, chass.PowerState, chass.Status.Health)
	}

	if len(args) != 0 && writer.RowCount() == 0 {
		return utils.ErrorExit(cmd, "chassis '%s' was not found.", args[0])
	}

	writer.Render()
	return nil
}
