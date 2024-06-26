// SPDX-License-Identifier: BSD-3-Clause
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage system connection information",
	Long: dedent.Dedent(`Manage system connection information.
	
	This command is used to save, update, and delete system connection
	information, as well as setting the default connection to use if
	not specified.`),
}

func init() {
	configCmd.AddCommand(NewGetConfigCmd())
	configCmd.AddCommand(NewAddConfigCmd())
	configCmd.AddCommand(NewRemoveConfigCmd())
	configCmd.AddCommand(NewSetConfigCmd())
	rootCmd.AddCommand(configCmd)
}

// NewGetConfigCmd returns a command for getting stpred system information.
func NewGetConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [CONNECTION_NAME]",
		Short: "Get saved connection information.",
		Long:  "The get command will get details for a specified connection or list all defined connections.",
		Run:   getConfigs,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

// getConfigs prints the saved system setting information.
func getConfigs(cmd *cobra.Command, args []string) {
	systems := config.GetSystems()
	defaultSys := config.GetDefault()

	writer := utils.NewTableWriter(cmd.OutOrStdout(), " ", "name", "user", "endpoint")
	if len(args) == 1 {
		system := config.GetSystem(args[0])
		isdefault := " "
		if args[0] == defaultSys {
			isdefault = "*"
		}

		writer.AddRow(
			isdefault, args[0], system.Username,
			fmt.Sprintf("%s://%s:%d", system.Protocol, system.Host, system.Port))
	} else {
		for system, settings := range systems {
			isdefault := " "
			if system == defaultSys {
				isdefault = "*"
			}

			writer.AddRow(
				isdefault, system, settings.Username,
				fmt.Sprintf("%s://%s:%d", settings.Protocol, settings.Host, settings.Port))
		}
	}

	writer.Render()
}

// NewAddConfigCmd creates a new subcommand for adding new system settings.
func NewAddConfigCmd() *cobra.Command {
	systemSettings := config.SystemConfig{}
	makeDefault := false
	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "Add new connection information.",
		Long:  "Adds new connection to use with the given name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addNewSystem(cmd, args[0], &systemSettings, makeDefault)
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().Uint16Var(&systemSettings.Port, "port", 0, "Port used to connect (defaults to 443, or port 80 if 'http' protocol is specified.)")
	cmd.Flags().StringVarP(&systemSettings.Host, "host", "e", "", "The host name or IP address of the system.")
	cmd.Flags().StringVarP(&systemSettings.Username, "user", "u", "", "The user name to connect as.")
	cmd.Flags().StringVarP(&systemSettings.Password, "password", "p", "", "The password to connect with.")
	cmd.Flags().StringVar(&systemSettings.Protocol, "protocol", "https", "Protocol to use (https (default) or http).")
	cmd.Flags().BoolVar(&systemSettings.Secure, "secure", false, "Enforce certificate validation with https connections (default allows self-signed certs).")
	cmd.Flags().BoolVar(&makeDefault, "default", false, "Set this connection as the default.")

	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")

	cmd.Flags().SortFlags = true

	return cmd
}

// addNewSystem performs the checks and handling for adding a new system.
func addNewSystem(cmd *cobra.Command, name string, settings *config.SystemConfig, makeDefault bool) error {
	system := config.GetSystem(name)
	if system != nil {
		return utils.ErrorExit(cmd, "a system named '%s' already exists. Delete and readd or update the existing one.\n", name)
	}

	// Validate the entries and set the conditional defaults
	if settings.Host == "" {
		settings.Host = name
	}

	if strings.Contains(settings.Host, "://") {
		// TODO: This is not complete, but I'm not going to worry about it for now
		// User provided actual endpoint, parse it out
		protoParts := strings.Split(settings.Host, "://")
		settings.Protocol = protoParts[0]
		if strings.Contains(protoParts[1], ":") {
			parts := strings.Split(protoParts[1], ":")
			settings.Host = parts[0]

			// Now make sure they didn't include the redfish path
			parts = strings.Split(parts[1], "/")
			port, err := strconv.Atoi(parts[0])
			if err != nil {
				return utils.ErrorExit(cmd, "failed to parse endpoint string: %s", err)
			}
			if port < 0 || port > 32768 {
				return utils.ErrorExit(cmd, "invalid port number provided.")
			}
			settings.Port = uint16(port)
		}
	}

	settings.Protocol = strings.ToLower(settings.Protocol)
	if settings.Port == 0 {
		if settings.Protocol == "http" {
			settings.Port = 80
		} else {
			settings.Port = 443
		}
	}

	// Add the new connection. We don't validate user name and password here. It
	// will be handled when they actually try to perform an operation.
	err := config.AddSystemConfig(name, settings, makeDefault)
	if err != nil {
		return utils.ErrorExit(cmd, "error adding system: %v", err)
	}
	getConfigs(cmd, []string{name})
	return nil
}

// NewRemoveConfigCmd returns a command for removing stored connection info.
func NewRemoveConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: "Remove stored connection information.",
		Long:  "Removes connection with the given name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.RemoveSystemConfig(args[0])
			if err != nil {
				return utils.ErrorExit(cmd, "error removing system: %v", err)
			}
			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

// NewSetConfigCmd returns a command for updating stored connection information.
func NewSetConfigCmd() *cobra.Command {
	systemSettings := config.SystemConfig{}
	makeDefault := false
	cmd := &cobra.Command{
		Use:   "set NAME",
		Short: "Set connection information.",
		Long:  "Updates stored connection with the given name with new values.",
		RunE: func(cmd *cobra.Command, args []string) error {
			system := config.GetSystem(args[0])
			if system == nil {
				return utils.ErrorExit(cmd, "connection '%s' not found, add new connection.", args[0])
			}

			defaultConnection := config.IsDefault(system)

			cmd.Flags().Visit(func(f *pflag.Flag) {
				switch f.Name {
				case "port":
					if systemSettings.Port > 32768 {
						fmt.Fprintf(cmd.ErrOrStderr(), "port value of %s is not valid.", f.Value.String())
						os.Exit(1)
					}

					system.Port = systemSettings.Port
				case "host":
					system.Host = systemSettings.Host
				case "user":
					system.Username = systemSettings.Username
				case "password":
					system.Password = systemSettings.Password
				case "protocol":
					system.Protocol = systemSettings.Protocol
				case "secure":
					system.Secure = systemSettings.Secure
				case "default":
					defaultConnection = makeDefault
				}
			})
			err := config.AddSystemConfig(args[0], system, defaultConnection)
			if err != nil {
				return utils.ErrorExit(cmd, "error adding system: %v", err)
			}
			getConfigs(cmd, []string{args[0]})
			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().Uint16VarP(&systemSettings.Port, "port", "p", 0, "Port used to connect (defaults to 443, or port 80 if 'http' protocol is specified.)")
	cmd.Flags().StringVarP(&systemSettings.Host, "host", "e", "", "The host name or IP address of the system.")
	cmd.Flags().StringVarP(&systemSettings.Username, "user", "u", "", "The user name to connect as.")
	cmd.Flags().StringVarP(&systemSettings.Password, "password", "s", "", "The password to connect with.")
	cmd.Flags().StringVar(&systemSettings.Protocol, "protocol", "https", "Protocol to use (https (default) or http).")
	cmd.Flags().BoolVar(&systemSettings.Secure, "secure", false, "Enforce certificate validation with https connections (default allows self-signed certs).")
	cmd.Flags().BoolVar(&makeDefault, "default", false, "Set this connection as the default.")

	cmd.Flags().SortFlags = true

	return cmd
}
