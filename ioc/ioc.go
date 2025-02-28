package ioc

import (
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/ffumaneri/github-cli/common/viper"
	github2 "github.com/ffumaneri/github-cli/github"
	ollama2 "github.com/ffumaneri/github-cli/ollama"
	"github.com/ffumaneri/github-cli/services"
	"github.com/google/go-github/v65/github"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"log"
	"net/url"
	"os"
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

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300   // size of the chunk is number of characters
	split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap

	ollamaWrapper := ollama2.NewOllamaWrapper(llm, &common.FS{}, func(llm2 ollama2.IOllamaLLM) embeddings.Embedder {
		ollamaEmbeder, err := embeddings.NewEmbedder(llm2.(*ollama.LLM))
		if err != nil {
			log.Fatal(err)
		}
		return ollamaEmbeder
	}, func(embeder embeddings.Embedder, collectionName string) vectorstores.VectorStore {
		config, err := common.NewConfig(viper.ViperLoadConfig)
		if err != nil {
			panic("error getting config")
		}
		// Create a new Qdrant vector store.
		quadrantUrl, err := url.Parse(config.QDrant_Url)
		if err != nil {
			log.Fatal(err)
		}

		store, err := qdrant.New(
			qdrant.WithURL(*quadrantUrl),
			qdrant.WithCollectionName(collectionName),
			qdrant.WithEmbedder(embeder),
		)
		if err != nil {
			log.Fatal(err)
		}
		return store
	}, func() textsplitter.TextSplitter {
		split := textsplitter.NewRecursiveCharacter()
		split.ChunkSize = 300   // size of the chunk is number of characters
		split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap
		return split
	}, func(filePath string, size int64) (documentloaders.Loader, *os.File) {
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file: ", err)
			return nil, nil
		}

		p := documentloaders.NewText(f)

		return p, f
	})
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
