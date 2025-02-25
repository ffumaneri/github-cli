package ollama

import (
	"context"
	"github.com/tmc/langchaingo/llms"
)

type IOllamaLLM interface {
	Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error)
}

func NewOllamaWrapper(llm IOllamaLLM) *OllamaWrapper {
	return &OllamaWrapper{llm}
}

type ILLMWrapper interface {
	AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error)
}
type OllamaWrapper struct {
	llm IOllamaLLM
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
