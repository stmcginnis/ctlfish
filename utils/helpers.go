//
// SPDX-License-Identifier: BSD-3-Clause
//
package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stmcginnis/ctlfish/config"
	"github.com/stmcginnis/gofish"
)

// ErrorExit is a helper function to format the error result of a command execution.
func ErrorExit(cmd *cobra.Command, message string, args ...interface{}) error {
	cmd.SilenceUsage = true
	return Error(message, args...)
}

// Error is a helper function to format the error results.
func Error(message string, args ...interface{}) error {
	msg := fmt.Sprintf(message, args...)
	return errors.New(msg)
}

// GofishClient will get a gofish client connection for the requested system.
// If connection == "", then the default system will be retrieved.
// The caller should close the client connection when done.
func GofishClient(connection string) (*gofish.APIClient, error) {
	// Create a new instance of gofish client
	var settings *config.SystemConfig
	if connection != "" {
		settings = config.GetSystem(connection)
	} else {
		settings = config.GetDefaultSystem()
	}

	if settings == nil {
		return nil, Error("unable to get system connection information.\nSet default to use or provide on command line with -c [NAME].")
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s:%d", settings.Protocol, settings.Host, settings.Port),
		Username: settings.Username,
		Password: settings.Password,
		Insecure: !settings.Secure,
	}

	c, err := gofish.Connect(config)
	if err != nil {
		return nil, Error("failed to connect to '%s': %v", config.Endpoint, err)
	}

	return c, nil
}

// BytesToReadable formats a byte count into a human readable representation.
func BytesToReadable(bytes int64) string {
	var val float32 = float32(bytes)

	if val < 1024 {
		return fmt.Sprintf("%0.2f Bytes", val)
	}

	val = val / 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f KB", val)
	}

	val = val / 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f MB", val)
	}

	val = val / 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f GB", val)
	}

	val = val / 1024
	return fmt.Sprintf("%0.2f TB", val)
}
