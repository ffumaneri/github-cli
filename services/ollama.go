package services

import (
	"context"
	"github.com/ffumaneri/github-cli/ollama"
)

type IOllamaService interface {
	AskLlm(prompt string) error
}

func NewOllamaService(llmWrapper ollama.ILLMWrapper, chunkConsumer func(chunk []byte)) *OllamaService {
	return &OllamaService{llmWrapper: llmWrapper, ChunkConsumer: chunkConsumer}
}

type OllamaService struct {
	ChunkConsumer func(chunk []byte)
	llmWrapper    ollama.ILLMWrapper
}

func (service *OllamaService) AskLlm(prompt string) (err error) {
	err = service.llmWrapper.AskLlm(prompt, func(ctx context.Context, chunk []byte) error {
		service.ChunkConsumer(chunk)
		return nil
	})
	return
}
