package common

import (
	"fmt"
)

// Config structure to hold the configuration
type Config struct {
	Token        string
	Owner        string
	Ollama_Model string
	QDrant_Url   string
}

var cachedConfig *Config // This will store the configuration as a singleton

type ConfigLoader func(path string, config interface{}) error

func NewConfig(configLoader ConfigLoader) (*Config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	config := &Config{}
	err := configLoader(".env", config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Ensure all required fields are set
	if config.Token == "" || config.Owner == "" || config.Ollama_Model == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	// More
	if config.QDrant_Url == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	cachedConfig = config
	return config, nil
}
