//
// SPDX-License-Identifier: BSD-3-Clause
//
package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/ctlfish/utils"
	"github.com/stmcginnis/gofish/redfish"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands for viewing and interacting with user accounts.",
}

func init() {
	userCmd.AddCommand(NewGetUserCmd())
	userCmd.AddCommand(NewSetUserCmd())
	rootCmd.AddCommand(userCmd)
	userCmd.PersistentFlags().StringP("connection", "c", config.GetDefault(), "The stored connection name to use.")
}

// NewGetUserCmd returns a command for getting user information.
func NewGetUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [NAME_OR_ID]",
		Short: "Get user information.",
		Long:  "Get details for a specified user or list all defined user.",
		RunE:  getUser,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

// getUser retrieves the user information from the system.
func getUser(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	as, err := c.Service.AccountService()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to access account service: %v", err)
	}

	users, err := as.Accounts()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve user information: %v", err)
	}

	writer := utils.NewTableWriter(cmd.OutOrStdout(), "name", "role", "enabled", "description")
	for _, user := range users {
		if len(args) == 1 && (user.ID != args[0] && user.Name != args[0] && user.UserName != args[0]) {
			continue
		}

		writer.AddRow(user.UserName, user.RoleID, user.Enabled, user.Description)
	}

	if len(args) != 0 && writer.RowCount() == 0 {
		return utils.ErrorExit(cmd, "user '%s' was not found.", args[0])
	}

	writer.Render()
	return nil
}

// NewSetUserCmd returns a command for updating user settings.
func NewSetUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set NAME",
		Short: "Set user settings.",
		Long:  "Updates password, role, or username.",
		RunE:  updateUser,
		Args:  cobra.ExactArgs(1),
	}

	cmd.Flags().StringP("username", "u", "", "New username for the user.")
	cmd.Flags().StringP("password", "p", "", "New password for the user.")
	cmd.Flags().StringP("role", "r", "", "The user role to apply.")
	cmd.Flags().SortFlags = true

	return cmd
}

// updateUser applies new settings to the user account.
func updateUser(cmd *cobra.Command, args []string) error {
	connection, _ := cmd.Flags().GetString("connection")
	c, err := utils.GofishClient(connection)
	if err != nil {
		return utils.ErrorExit(cmd, err.Error())
	}
	defer c.Logout()

	as, err := c.Service.AccountService()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to access account service: %v", err)
	}

	users, err := as.Accounts()
	if err != nil {
		return utils.ErrorExit(cmd, "failed to retrieve user information: %v", err)
	}

	var user *redfish.ManagerAccount
	user = nil

	for _, account := range users {
		if account.ID != args[0] && account.Name != args[0] && account.UserName != args[0] {
			continue
		}
		user = account
		break
	}

	usernameFlag := cmd.Flag("username")
	if usernameFlag.Changed {
		// TODO: since we retrieved all accounts, might be good to add validation
		// here that the newly requested username does not conflict with another
		// account.
		user.UserName = usernameFlag.Value.String()
	}

	passwordFlag := cmd.Flag("password")
	if passwordFlag.Changed {
		// Make sure it meetings the criteria
		minPassLen := as.MinPasswordLength
		maxPassLen := as.MaxPasswordLength

		newPass := passwordFlag.Value.String()
		if len(newPass) < minPassLen || len(newPass) > maxPassLen {
			return utils.ErrorExit(cmd, "account password must be between %d - %d in length", minPassLen, maxPassLen)
		}

		user.Password = newPass
	}

	roleFlag := cmd.Flag("role")
	if roleFlag.Changed {
		// Validate the role being set
		roles, err := as.Roles()
		if err != nil {
			return utils.ErrorExit(cmd, "unable to retrieve available roles: %s", err.Error())
		}

		newRole := strings.ToLower(roleFlag.Value.String())
		roleFound := false
		for _, role := range roles {
			if strings.EqualFold(newRole, role.Name) || strings.EqualFold(newRole, role.ID) {
				user.RoleID = role.ID
				roleFound = true
				break
			}
		}

		if !roleFound {
			return utils.ErrorExit(cmd, "role '%s' was not found on this system", newRole)
		}
	}

	err = user.Update()
	if err != nil {
		return utils.ErrorExit(cmd, "error updating user '%s': %s", user.UserName, err.Error())
	}

	writer := utils.NewTableWriter(cmd.OutOrStdout(), "name", "role", "enabled", "description")
	writer.AddRow(user.UserName, user.RoleID, user.Enabled, user.Description)
	writer.Render()
	return nil
}
