// SPDX-License-Identifier: BSD-3-Clause
package set

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/gofish/redfish"

	"github.com/stmcginnis/ctlfish/utils"
)

var userCmd = &cobra.Command{
	Use:     "user [NAME_OR_ID]",
	Aliases: []string{"u"},
	Short:   "Set user information.",
	Long:    "Updates password, role, or username.",
	RunE:    updateUser,
	Args:    cobra.ExactArgs(1),
}

func init() {
	userCmd.Flags().StringP("username", "u", "", "New username for the user.")
	userCmd.Flags().StringP("password", "p", "", "New password for the user.")
	userCmd.Flags().StringP("role", "r", "", "The user role to apply.")
	userCmd.Flags().SortFlags = true
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
