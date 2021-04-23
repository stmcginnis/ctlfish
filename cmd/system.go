//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

// systemCmd represents the system command
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for viewing and interacting with system objects",
}

func init() {
	systemCmd.AddCommand(NewGetSystemCmd())
	systemCmd.AddCommand(NewResetSystemCmd())
	rootCmd.AddCommand(systemCmd)
	systemCmd.PersistentFlags().StringP("connection", "c", config.GetDefault(), "The stored connection name to use.")
}

// NewGetSystemCmd returns a command for getting system information.
func NewGetSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [NAME_OR_ID]",
		Short: "Get system information.",
		RunE:  getSystem,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
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

// NewResetSystemCmd returns a command for resetting a system.
func NewResetSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset [NAME_OR_ID]",
		Short: "Reset a system.",
		RunE:  resetSystem,
		Args:  cobra.ExactArgs(1),
	}

	return cmd
}

// resetSystem performs a reset on a given system.
func resetSystem(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	var sys *redfish.ComputerSystem
	systems, err := c.Service.Systems()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve system information: %v", err)
	}

	for _, system := range systems {
		if system.Name == args[0] || system.ID == args[0] {
			sys = system
			break
		}
	}

	if sys == nil {
		return utils.ErrorExit(cmd, "unable to locate system '%s'", args[0])
	}

	// There are different types of resets that can be performed. We may want to
	// support letting the user specify, but for now just default to PowerCycle.
	err = sys.Reset(redfish.PowerCycleResetType)
	if err != nil {
		msg := err.Error()
		if rfErr, ok := err.(*common.Error); ok {
			msg = rfErr.Message
		}
		return utils.ErrorExit(cmd, "error performing reset: %v", msg)
	}
	return nil
}
