// SPDX-License-Identifier: BSD-3-Clause
package get

import (
	"github.com/spf13/cobra"

	"github.com/stmcginnis/ctlfish/utils"
)

var userCmd = &cobra.Command{
	Use:     "user [NAME_OR_ID]",
	Aliases: []string{"users", "u"},
	Short:   "Get user information.",
	Long:    "Get details for a specified user or list all defined users.",
	RunE:    getUser,
	Args:    cobra.MaximumNArgs(1),
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
