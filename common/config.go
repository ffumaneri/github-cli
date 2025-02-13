package common

import (
	"fmt"
	"github.com/google/go-github/v65/github"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

// Config structure to hold the configuration
type Config struct {
	Token        string
	Owner        string
	Ollama_Model string
}

var cachedConfig *Config // This will store the configuration as a singleton

func GetConfig() (*Config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Ensure all required fields are set
	if config.Token == "" || config.Owner == "" || config.Ollama_Model == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	cachedConfig = config
	return config, nil
}

func GetGithubClient() (*github.Client, string, error) {
	config, err := GetConfig()
	if err != nil {
		panic("error getting config")
	}
	// Create Git hub client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}

func GetOllamaClient() (llm *ollama.LLM, ok bool) {
	config, err := GetConfig()
	if err != nil {
		return nil, false
	}
	llm, err = ollama.New(ollama.WithModel(config.Ollama_Model))
	if err != nil {
		log.Fatal(err)
	}
	return llm, true
}
