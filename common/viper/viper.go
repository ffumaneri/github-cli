package viper

import (
	"fmt"
	"github.com/spf13/viper"
)

// LoadConfig is a wrapper function to load configuration from the given file.
func ViperLoadConfig(path string, config interface{}) error {
	// Set the file to read
	viper.SetConfigFile(path)

	// Read in the config file
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the config into the provided struct
	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	return nil
}
