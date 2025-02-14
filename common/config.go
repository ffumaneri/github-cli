package common

import (
	"context"
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

type configLoader func(path string, config interface{}) error

// LoadConfig is a wrapper function to load configuration from the given file.
func viperLoadConfig(path string, config interface{}) error {
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

func newConfig(configLoader configLoader) (*Config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	config := &Config{}
	err := configLoader(".env", config)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Ensure all required fields are set
	if config.Token == "" || config.Owner == "" || config.Ollama_Model == "" {
		return nil, fmt.Errorf("missing required configuration values")
	}

	cachedConfig = config
	return config, nil
}

func GetGithubContext() context.Context {
	config, err := newConfig(viperLoadConfig)
	if err != nil {
		panic("error getting config")
	}
	ghClient, owner, err := getGithubClient(config)
	if err != nil {
		panic("error getting config")
	}
	ctx := context.WithValue(context.TODO(), "githubClient", ghClient)
	ctx = context.WithValue(ctx, "githubOwner", owner)

	llmClient, ok := getOllamaClient(config)
	if !ok {
		panic("error getting config")
	}
	ctx = context.WithValue(ctx, "ollamaClient", llmClient)
	return ctx
}

func GetLLMContext() context.Context {
	config, err := newConfig(viperLoadConfig)
	if err != nil {
		panic("error getting config")
	}

	llmClient, ok := getOllamaClient(config)
	if !ok {
		panic("error getting config")
	}
	ctx := context.WithValue(context.TODO(), "ollamaClient", llmClient)
	return ctx
}

func getGithubClient(config *Config) (*github.Client, string, error) {
	// Create Github client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}

func getOllamaClient(config *Config) (llm *ollama.LLM, ok bool) {
	llm, err := ollama.New(ollama.WithModel(config.Ollama_Model))
	if err != nil {
		panic(fmt.Errorf("error getting ollama client with model %s: %w", config.Ollama_Model, err))
	}
	return llm, true
}
