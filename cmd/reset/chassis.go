// SPDX-License-Identifier: BSD-3-Clause
package reset

import (
	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"

	"github.com/stmcginnis/ctlfish/utils"
)

var chassisCmd = &cobra.Command{
	Use:     "chassis [NAME_OR_ID]",
	Aliases: []string{"c"},
	Short:   "Reset a chassis.",
	RunE:    resetChassis,
	Args:    cobra.ExactArgs(1),
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
