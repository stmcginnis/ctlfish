//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
	"github.com/stmcginnis/gofish/redfish"
)

// driveCmd represents the drive command
var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Commands for viewing and interacting with drive objects",
}

func init() {
	driveCmd.AddCommand(NewGetDriveCmd())
	rootCmd.AddCommand(driveCmd)
	driveCmd.PersistentFlags().StringP("connection", "c", config.GetDefault(), "The stored connection name to use.")
}

// NewGetDriveCmd returns a command for getting drive information.
func NewGetDriveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [NAME_OR_ID]",
		Short: "Get drive information.",
		RunE:  getDrive,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

// getDrive retrieves the drive information.
func getDrive(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	chassiss, err := c.Service.Chassis()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve system information: %v", err)
	}

	// TODO: we may need to look between the System, Chassis, and Storage root
	// objects to collect all drive information. Also need to check if we need
	// to scrub in case there are multiple relationships to the same drives.

	// Collect drive info from all chassis
	drives := []*redfish.Drive{}
	for _, chassis := range chassiss {
		chassisDrives, err := chassis.Drives()
		if err != nil {
			return utils.ErrorExit(cmd, "failed to get drive information for chassis %q", chassis.Name)
		}

		drives = append(drives, chassisDrives...)
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
