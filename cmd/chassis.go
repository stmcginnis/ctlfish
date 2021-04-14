//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
	"github.com/stmcginnis/gofish"
)

// chassisCmd represents the chassis command
var chassisCmd = &cobra.Command{
	Use:   "chassis",
	Short: "Commands for viewing and interacting with chassis objects.",
}

func init() {
	chassisCmd.AddCommand(NewGetChassisCmd())
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
	// Create a new instance of gofish client
	var settings *config.SystemConfig
	connection, _ := cmd.Flags().GetString("connection")
	if connection != "" {
		settings = config.GetSystem(connection)
	} else {
		settings = config.GetDefaultSystem()
	}

	if settings == nil {
		return utils.ErrorExit(cmd, "unable to get system connection information.\nSet default to use or provide on command line with -c [NAME].")
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s:%d", settings.Protocol, settings.Host, settings.Port),
		Username: settings.Username,
		Password: settings.Password,
		Insecure: !settings.Secure,
	}

	c, err := gofish.Connect(config)
	if err != nil {
		return utils.ErrorExit(cmd, "failed to connect to '%s': %v", config.Endpoint, err)
	}
	defer c.Logout()

	chassis, err := c.Service.Chassis()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve chassis information: %v", err)
	}

	var data [][]string
	for _, chass := range chassis {
		if len(args) == 1 && (chass.ID != args[0] && chass.Name != args[0]) {
			continue
		}

		row := []string{
			chass.ID,
			chass.Name,
			string(chass.PowerState),
			string(chass.Status.Health),
		}
		data = append(data, row)
	}

	if len(args) != 0 && len(data) == 0 {
		return utils.ErrorExit(cmd, "chassis '%s' was not found.", args[0])
	}

	headers := []string{"id", "name", "power", "status"}
	utils.PrintTable(headers, data)
	return nil
}
