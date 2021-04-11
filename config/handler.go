//
// SPDX-License-Identifier: BSD-3-Clause
//
package config

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var appConfig Config

// SystemConfig is the config settings for our systems.
type SystemConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Protocol string `yaml:"protocol"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Secure   bool   `yaml:"secure"`
}

// Config is the configuration settings we use.
type Config struct {
	Default string                  `yaml:"default"`
	Systems map[string]SystemConfig `yaml:"systems"`
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ctlfish" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ctlfish")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("CTLFISH")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetDefault("default", "")
	viper.SetDefault("systems", (&Config{}).Systems)

	// Write out config file so it is created on first run
	_ = viper.SafeWriteConfig()

	loadConfig()
}

func loadConfig() {
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}
}

// GetSystems gets all configured system settings.
func GetSystems() map[string]SystemConfig {
	return appConfig.Systems
}

// GetSystem gets the config settings for the system with the given host name.
func GetSystem(name string) *SystemConfig {
	for systemName, system := range appConfig.Systems {
		if systemName == name {
			return &system
		}
	}

	// Not found by the name, let's see if they provided the hostname
	for _, system := range appConfig.Systems {
		if system.Host == name {
			return &system
		}
	}

	return nil
}

// GetDefaultSystem gets the config settings for the system set as the default.
func GetDefaultSystem() *SystemConfig {
	if appConfig.Default == "" && len(appConfig.Systems) == 1 {
		// Only one system anyway, return that
		for _, system := range appConfig.Systems {
			return &system
		}
	}

	for systemName, system := range appConfig.Systems {
		if systemName == appConfig.Default {
			return &system
		}
	}

	return nil
}

// IsDefault checks if a system is the current default
func IsDefault(system *SystemConfig) bool {
	defaultSystem := GetDefaultSystem()
	if defaultSystem == nil {
		return false
	}

	return (system.Host == defaultSystem.Host &&
		system.Username == defaultSystem.Username &&
		system.Password == defaultSystem.Password)
}

// AddSystemConfig adds or updates system config settings. If system already
// exists it will be overwritten with the new settings.
func AddSystemConfig(name string, sysConfig *SystemConfig, makeDefault bool) {
	appConfig.Systems[name] = *sysConfig
	if makeDefault {
		appConfig.Default = name
		viper.Set("default", name)
	}

	viper.Set("systems", appConfig.Systems)
	viper.WriteConfig()
}

// RemoveSystemConfig removes system config settings. If the system being removed
// was the default connection, default is set to nothing.
func RemoveSystemConfig(name string) {
	if appConfig.Default == name {
		appConfig.Default = ""
	}

	delete(appConfig.Systems, name)
	viper.Set("systems", appConfig.Systems)
	viper.WriteConfig()
}

// SetDefault will set the default system to use if not explicitly provided.
func SetDefault(name string) error {
	// Validate that we actually have a system named this
	for system := range appConfig.Systems {
		if system == name {
			appConfig.Default = name
			viper.Set("default", name)
			viper.WriteConfig()
			return nil
		}
	}

	return fmt.Errorf("no system named %s", name)
}

// GetDefault returns the name of the system to be used as the default connection.
func GetDefault() string {
	return appConfig.Default
}
