package ollama

import (
	"context"
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func NewOllamaClient(config *common.Config) (llm *ollama.LLM, ok bool) {
	llm, err := ollama.New(ollama.WithModel(config.Ollama_Model))
	if err != nil {
		panic(fmt.Errorf("error getting ollama client with model %s: %w", config.Ollama_Model, err))
	}
	return llm, true
}

func NewOllamaWrapper(llm *ollama.LLM) *OllamaWrapper {
	return &OllamaWrapper{llm}
}

type ILLMWrapper interface {
	AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error)
}
type OllamaWrapper struct {
	llm *ollama.LLM
}

func (o *OllamaWrapper) AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error) {

	llmCtx := context.Background()
	completion, err := o.llm.Call(llmCtx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(streamingFunc),
	)
	_ = completion
	return
}
