package common

import (
	"errors"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name             string
		mockConfigLoader ConfigLoader
		expectedError    bool
	}{
		{
			name: "Valid configuration",
			mockConfigLoader: func(path string, config interface{}) error {
				cfg, ok := config.(*Config)
				if !ok {
					return errors.New("invalid config type")
				}
				cfg.Token = "validToken"
				cfg.Owner = "validOwner"
				cfg.Ollama_Model = "validModel"
				return nil
			},
			expectedError: false,
		},
		{
			name: "Empty Token",
			mockConfigLoader: func(path string, config interface{}) error {
				cfg, ok := config.(*Config)
				if !ok {
					return errors.New("invalid config type")
				}
				cfg.Token = ""
				cfg.Owner = "validOwner"
				cfg.Ollama_Model = "validModel"
				return nil
			},
			expectedError: true,
		},
		{
			name: "Empty Owner",
			mockConfigLoader: func(path string, config interface{}) error {
				cfg, ok := config.(*Config)
				if !ok {
					return errors.New("invalid config type")
				}
				cfg.Token = "validToken"
				cfg.Owner = ""
				cfg.Ollama_Model = "validModel"
				return nil
			},
			expectedError: true,
		},
		{
			name: "Empty Ollama_Model",
			mockConfigLoader: func(path string, config interface{}) error {
				cfg, ok := config.(*Config)
				if !ok {
					return errors.New("invalid config type")
				}
				cfg.Token = "validToken"
				cfg.Owner = "validOwner"
				cfg.Ollama_Model = ""
				return nil
			},
			expectedError: true,
		},
		{
			name: "ConfigLoader error",
			mockConfigLoader: func(path string, config interface{}) error {
				return errors.New("failed to load configuration")
			},
			expectedError: true,
		},
	}

	// Reset cachedConfig for testing
	cachedConfig = nil

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cachedConfig = nil // Ensure cachedConfig is nil for each test
			_, err := NewConfig(tt.mockConfigLoader)
			if (err != nil) != tt.expectedError {
				t.Errorf("NewConfig() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
