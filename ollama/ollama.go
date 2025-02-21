package ollama

import (
	"context"
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

func NewOllamaClient(config *common.Config) (llm *ollama.LLM, ok bool) {
	llm, err := ollama.New(ollama.WithModel(config.Ollama_Model))
	if err != nil {
		panic(fmt.Errorf("error getting ollama client with model %s: %w", config.Ollama_Model, err))
	}
	return llm, true
}

func AskLlm(llm *ollama.LLM, prompt string) (err error) {

	llmCtx := context.Background()
	completion, err := llm.Call(llmCtx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = completion

	return
}
