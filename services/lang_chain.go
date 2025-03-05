package services

import (
	"context"
	"github.com/ffumaneri/github-cli/lang_chain"
)

type ILangChainService interface {
	AskLlm(prompt string) error
	LoadSourceCode(name, directory string) error
	AskLlmWithContext(contextName, prompt string) (err error)
}

func NewLangChainService(llmWrapper lang_chain.ILangChainWrapper, chunkConsumer func(chunk []byte)) *LangChainService {
	return &LangChainService{llmWrapper: llmWrapper, ChunkConsumer: chunkConsumer}
}

type LangChainService struct {
	ChunkConsumer func(chunk []byte)
	llmWrapper    lang_chain.ILangChainWrapper
}

func (service *LangChainService) AskLlm(prompt string) (err error) {
	err = service.llmWrapper.AskLlm(prompt, func(ctx context.Context, chunk []byte) error {
		service.ChunkConsumer(chunk)
		return nil
	})
	return
}
func (service *LangChainService) AskLlmWithContext(contextName, prompt string) (err error) {
	err = service.llmWrapper.AskWithContext(contextName, prompt)
	return
}
func (service *LangChainService) LoadSourceCode(name, directory string) (err error) {
	err = service.llmWrapper.LoadSourceCode(name, directory)
	return
}
