package services

import (
	"context"
	ollama2 "github.com/ffumaneri/github-cli/ollama"
	"github.com/tmc/langchaingo/llms/ollama"
)

type IOllamaService interface {
	AskLlm(prompt string) error
}

type OllamaService struct {
	llm           *ollama.LLM
	ChunkConsumer func(chunk []byte)
}

func NewOllamaService(llm *ollama.LLM, chunkConsumer func(chunk []byte)) *OllamaService {
	return &OllamaService{llm: llm, ChunkConsumer: chunkConsumer}
}

func (service *OllamaService) AskLlm(prompt string) (err error) {
	err = ollama2.AskLlm(service.llm, prompt, func(ctx context.Context, chunk []byte) error {
		service.ChunkConsumer(chunk)
		return nil
	})
	return
}
