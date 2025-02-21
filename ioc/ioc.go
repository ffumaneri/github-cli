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

func GetGithubClient() (*github.Client, string) {
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

func NewGithubService() *services.GithubService {
	ghClient, owner := GetGithubClient()
	return services.NewGithubService(ghClient, owner)
}

func GetLLMClient() (llm *ollama.LLM) {
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

func NewOllamaService() *services.OllamaService {
	llm := GetLLMClient()
	return services.NewOllamaService(llm)
}
