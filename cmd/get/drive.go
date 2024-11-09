// SPDX-License-Identifier: BSD-3-Clause
package get

import (
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"

	"github.com/stmcginnis/ctlfish/utils"
)

var driveCmd = &cobra.Command{
	Use:     "drive [NAME_OR_ID]",
	Aliases: []string{"drives", "d"},
	Short:   "Get drive information.",
	RunE:    getDrive,
	Args:    cobra.MaximumNArgs(1),
}

// getDrive retrieves the drive information.
func getDrive(cmd *cobra.Command, args []string) error {
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

	// Collect drive info from all systems
	drives := []*redfish.Drive{}
	for _, sys := range systems {
		storage, err := sys.Storage()
		if err != nil {
			return utils.ErrorExit(
				cmd,
				"failed to get drive information for system %q",
				sys.Name)
		}

		for _, stor := range storage {
			storageDrives, err := stor.Drives()
			if err != nil {
				// Some storage systems do not contain drives
				continue
			}

			drives = append(drives, storageDrives...)
		}
	}

	writer := utils.NewTableWriter(
		cmd.OutOrStdout(),
		"name", "size", "status", "manufacturer", "model", "serial number")
	for _, drive := range drives {
		if len(args) == 1 && (drive.ID != args[0] && drive.Name != args[0]) {
			continue
		}

		writer.AddRow(
			drive.Name,
			utils.BytesToReadable(drive.CapacityBytes),
			drive.Status.Health,
			drive.Manufacturer,
			drive.Model,
			drive.SerialNumber)
	}

	if len(args) != 0 && writer.RowCount() == 0 {
		return utils.ErrorExit(cmd, "drive '%s' was not found.", args[0])
	}

	writer.Render()
	return nil
}
