package services

import (
	ollama2 "github.com/ffumaneri/github-cli/ollama"
	"github.com/tmc/langchaingo/llms/ollama"
)

type OllamaService struct {
	llm *ollama.LLM
}

func NewOllamaService(llm *ollama.LLM) *OllamaService {
	return &OllamaService{llm: llm}
}

func (service *OllamaService) AskLlm(prompt string) (err error) {
	err = ollama2.AskLlm(service.llm, prompt)
	if err != nil {
		return err
	}
	return
}
