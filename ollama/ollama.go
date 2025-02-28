package ollama

import (
	"context"
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"log"
	"os"
)

type IOllamaLLM interface {
	Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error)
}

func NewOllamaWrapper(llm IOllamaLLM, fs common.IFS, embedder func(llm IOllamaLLM) embeddings.Embedder, store func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore, splitter func() textsplitter.TextSplitter, documentLoader func(filePath string, size int64) (documentloaders.Loader, *os.File)) *OllamaWrapper {
	return &OllamaWrapper{llm, fs, embedder, store, splitter, documentLoader}
}

type ILLMWrapper interface {
	AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error)
	LoadSourceCode(name, directory string) (err error)
}
type OllamaWrapper struct {
	llm            IOllamaLLM
	fs             common.IFS
	embedder       func(llm IOllamaLLM) embeddings.Embedder
	store          func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore
	splitter       func() textsplitter.TextSplitter
	documentLoader func(filePath string, size int64) (documentloaders.Loader, *os.File)
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

func (o *OllamaWrapper) LoadSourceCode(name, directory string) (err error) {
	err = o.fs.WalkDir(directory, func(path string, size int64) {
		docs, err := o.processDocument(path, size)
		if err != nil {
			log.Fatal(err)
		}
		store := o.store(o.embedder(o.llm), name)
		_, err = store.AddDocuments(context.Background(), docs)
		if err != nil {
			log.Fatal(err)
		}
	})
	return
}

func (o *OllamaWrapper) processDocument(path string, size int64) ([]schema.Document, error) {
	split := o.splitter()
	p, f := o.documentLoader(path, size)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Error closing file: ", err)
			return
		}
	}(f)

	docs, err := p.LoadAndSplit(context.Background(), split)
	if err != nil {
		fmt.Println("Error loading document: ", err)
		return nil, err
	}
	log.Println("Documents loaded: ", len(docs))
	return docs, nil
}
