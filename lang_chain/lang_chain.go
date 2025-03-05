package lang_chain

import (
	"context"
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"log"
	"os"
)

func NewLangChainWrapper(llm llms.Model, fs common.IFS, embedder func(llm llms.Model) embeddings.Embedder, store func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore, splitter func() textsplitter.TextSplitter, documentLoader func(filePath string, size int64) (documentloaders.Loader, *os.File), retrievalQA func(model llms.Model, retriever vectorstores.Retriever) chains.Chain) *LangChainWrapper {
	return &LangChainWrapper{llm, fs, embedder, store, splitter, documentLoader, retrievalQA}
}

type ILangChainWrapper interface {
	AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error)
	LoadSourceCode(name, directory string) (err error)
	AskWithContext(contextName, searchQuery string) (err error)
}
type LangChainWrapper struct {
	llm            llms.Model
	fs             common.IFS
	embedder       func(llm llms.Model) embeddings.Embedder
	store          func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore
	splitter       func() textsplitter.TextSplitter
	documentLoader func(filePath string, size int64) (documentloaders.Loader, *os.File)
	retrievalQA    func(llm llms.Model, retriever vectorstores.Retriever) chains.Chain
}

func (o *LangChainWrapper) AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) (err error) {

	llmCtx := context.Background()
	completion, err := o.llm.Call(llmCtx, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(streamingFunc),
	)
	_ = completion
	return
}

func (o *LangChainWrapper) LoadSourceCode(name, directory string) (err error) {
	err = o.fs.WalkDir(directory, func(path string, size int64) {
		docs, err := o.processDocument(path, size)
		if err != nil {
			log.Fatal(err)
		}
		e := o.embedder(o.llm)
		store := o.store(e, name)
		_, err = store.AddDocuments(context.Background(), docs)
		if err != nil {
			log.Fatal(err)
		}
	})
	return
}
func (o *LangChainWrapper) AskWithContext(contextName, searchQuery string) (err error) {
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		//vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	store := o.store(o.embedder(o.llm), contextName)
	retriever := vectorstores.ToRetriever(store, 10, optionsVector...)
	//search
	//resDocs, err := retriever.GetRelevantDocuments(context.Background(), searchQuery)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for i, doc := range resDocs {
	//	fmt.Printf("Doc: %d Score: %f, content: %s metadata: %s\n", i, doc.Score, doc.PageContent)
	//}

	retrievalQA := o.retrievalQA(o.llm, retriever)
	var values map[string]any = make(map[string]any)
	values["query"] = searchQuery
	call, err := retrievalQA.Call(context.Background(), values)
	if err != nil {
		return
	}
	for key, values := range call {
		fmt.Printf("Key: %s, Values: %v\n", key, values)
	}
	return
}

func (o *LangChainWrapper) processDocument(path string, size int64) ([]schema.Document, error) {
	split := o.splitter()
	p, f := o.documentLoader(path, size)
	defer func(f *os.File) {
		if f == nil {
			return
		}
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
