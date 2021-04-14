//
// SPDX-License-Identifier: BSD-3-Clause
//
package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// ErrorExit is a helper function to format the error result of a command execution.
func ErrorExit(cmd *cobra.Command, message string, args ...interface{}) error {
	cmd.SilenceUsage = true
	msg := fmt.Sprintf(message, args...)
	return errors.New(msg)
}
