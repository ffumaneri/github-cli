package ioc

import (
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/ffumaneri/github-cli/common/viper"
	github2 "github.com/ffumaneri/github-cli/github"
	ollama2 "github.com/ffumaneri/github-cli/ollama"
	"github.com/ffumaneri/github-cli/services"
	"github.com/google/go-github/v65/github"
	"github.com/tmc/langchaingo/llms/ollama"
)

// Container defines an interface for initializing services and clients.
type Container interface {
	NewGithubService() services.IGithubService
	NewOllamaService() services.IOllamaService
}

// AppContainer is a concrete implementation of Container.
type AppContainer struct{}

func (ioc *AppContainer) NewGithubService() services.IGithubService {
	ghClient, owner := ioc.getGithubClient()
	ghWrapper := github2.NewGithubWrapper(ghClient, owner)
	return services.NewGithubService(owner, ghWrapper, func(data string) {
		fmt.Println(data)
	})
}

func (ioc *AppContainer) NewOllamaService() services.IOllamaService {
	llm := ioc.getLLMClient()
	ollamaWrapper := ollama2.NewOllamaWrapper(llm)
	return services.NewOllamaService(ollamaWrapper, func(chunk []byte) {
		fmt.Print(string(chunk))
	})
}
func NewGithubClient(config *common.Config) (*github.Client, string, error) {
	// Create Github client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}

func (ioc *AppContainer) getGithubClient() (*github.Client, string) {
	config, err := common.NewConfig(viper.ViperLoadConfig)
	if err != nil {
		panic("error getting config")
	}
	ghClient, owner, err := NewGithubClient(config)
	if err != nil {
		panic("error getting config")
	}
	return ghClient, owner
}
func NewOllamaClient(config *common.Config) (llm *ollama.LLM, ok bool) {
	llm, err := ollama.New(ollama.WithModel(config.Ollama_Model))
	if err != nil {
		panic(fmt.Errorf("error getting ollama client with model %s: %w", config.Ollama_Model, err))
	}
	return llm, true
}
func (ioc *AppContainer) getLLMClient() *ollama.LLM {
	config, err := common.NewConfig(viper.ViperLoadConfig)
	if err != nil {
		panic("error getting config")
	}
	llmClient, ok := NewOllamaClient(config)
	if !ok {
		panic("error getting config")
	}
	return llmClient
}
