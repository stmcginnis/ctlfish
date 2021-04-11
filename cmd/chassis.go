//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/printers"
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
		Run:   getChassis,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

// getChassis retrieves the chassis information from the system.
func getChassis(cmd *cobra.Command, args []string) {
	// Create a new instance of gofish client
	var settings *config.SystemConfig
	connection, _ := cmd.Flags().GetString("connection")
	if connection != "" {
		settings = config.GetSystem(connection)
	} else {
		settings = config.GetDefaultSystem()
	}

	if settings == nil {
		fmt.Fprintln(os.Stderr, "Unable to get system connection information.")
		fmt.Fprintln(os.Stderr, "Set default to use or provide on command line with -c [NAME].")
		os.Exit(1)
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s:%d", settings.Protocol, settings.Host, settings.Port),
		Username: settings.Username,
		Password: settings.Password,
		Insecure: !settings.Secure,
	}

	c, err := gofish.Connect(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to '%s': %v\n", config.Endpoint, err)
		os.Exit(1)
	}
	defer c.Logout()

	chassis, err := c.Service.Chassis()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve chassis information: %v\n", err)
		os.Exit(1)
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

	headers := []string{"id", "name", "power", "status"}
	printers.PrintTable(headers, data)
}
