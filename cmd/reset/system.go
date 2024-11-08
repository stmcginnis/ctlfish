// SPDX-License-Identifier: BSD-3-Clause
package reset

import (
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"

	"github.com/stmcginnis/ctlfish/utils"
)

var systemCmd = &cobra.Command{
	Use:     "system [NAME_OR_ID]",
	Aliases: []string{"s"},
	Short:   "Reset a system.",
	RunE:    resetSystem,
	Args:    cobra.ExactArgs(1),
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
