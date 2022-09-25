// Package configreader is responsible for getting needed configs.
package configreader

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// SetConfigFile defines path and name of the desired config file.
func SetConfigFile(path string) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config read: %w", err)
	}

	return nil
}

// GetString reads string with the specified key from the config file declared in SetConfigFile.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetDuration reads time.Duration with the specified key from the config file declared in SetConfigFile.
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetInt reads int with the specified key from the config file declared in SetConfigFile.
func GetInt(key string) int {
	return viper.GetInt(key)
}
