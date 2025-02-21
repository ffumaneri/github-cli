package ioc

import (
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
	return services.NewGithubService(ghClient, owner)
}

func (ioc *AppContainer) NewOllamaService() services.IOllamaService {
	llm := ioc.getLLMClient()
	return services.NewOllamaService(llm)
}

func (ioc *AppContainer) getGithubClient() (*github.Client, string) {
	config, err := common.NewConfig(viper.ViperLoadConfig)
	if err != nil {
		panic("error getting config")
	}
	ghClient, owner, err := github2.NewGithubClient(config)
	if err != nil {
		panic("error getting config")
	}
	return ghClient, owner
}

func (ioc *AppContainer) getLLMClient() *ollama.LLM {
	config, err := common.NewConfig(viper.ViperLoadConfig)
	if err != nil {
		panic("error getting config")
	}
	llmClient, ok := ollama2.NewOllamaClient(config)
	if !ok {
		panic("error getting config")
	}
	return llmClient
}
