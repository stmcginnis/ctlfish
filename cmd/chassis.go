//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

// chassisCmd represents the chassis command
var chassisCmd = &cobra.Command{
	Use:   "chassis",
	Short: "Commands for viewing and interacting with chassis objects.",
}

func init() {
	chassisCmd.AddCommand(NewGetChassisCmd())
	chassisCmd.AddCommand(NewResetChassisCmd())
	rootCmd.AddCommand(chassisCmd)
	chassisCmd.PersistentFlags().StringP("connection", "c", config.GetDefault(), "The stored connection name to use.")
}

// NewGetChassisCmd returns a command for getting chassis information.
func NewGetChassisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [NAME_OR_ID]",
		Short: "Get chassis information.",
		Long:  "Get details for a specified chassis or list all defined chassis.",
		RunE:  getChassis,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
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

// NewResetChassisCmd returns a command for resetting a chassis.
func NewResetChassisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset [NAME_OR_ID]",
		Short: "Reset a chassis.",
		RunE:  resetChassis,
		Args:  cobra.ExactArgs(1),
	}

	return cmd
}

// resetChassis performs a reset on a given chassis.
func resetChassis(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	var ch *redfish.Chassis
	chassis, err := c.Service.Chassis()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve chassis information: %v", err)
	}

	for _, chass := range chassis {
		if chass.Name == args[0] || chass.ID == args[0] {
			ch = chass
			break
		}
	}

	if ch == nil {
		return utils.ErrorExit(cmd, "unable to locate chassis '%s'", args[0])
	}

	// There are different types of resets that can be performed. We may want to
	// support letting the user specify, but for now just default to PowerCycle.
	err = ch.Reset(redfish.PowerCycleResetType)
	if err != nil {
		msg := err.Error()
		if rfErr, ok := err.(*common.Error); ok {
			msg = rfErr.Message
		}
		return utils.ErrorExit(cmd, "error performing reset: %v", msg)
	}
	return nil
}
