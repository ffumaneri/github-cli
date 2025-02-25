package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"log"
	"net/url"
	"os"
)

const URL = "http://localhost:6333"
const COLLECTION_NAME = "ioc_go"
const MODEL = "llama3.2"

func Search() {
	llm, err := ollama.New(ollama.WithModel(MODEL))

	if err != nil {
		log.Fatal(err)
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}

	quadrantUrl, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	store, err := qdrant.New(
		qdrant.WithURL(*quadrantUrl),
		qdrant.WithCollectionName(COLLECTION_NAME),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		log.Fatal(err)
	}
	searchQuery := "how does it create a OllamaService?"
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		//vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	retriever := vectorstores.ToRetriever(store, 10, optionsVector...)
	// search
	resDocs, err := retriever.GetRelevantDocuments(context.Background(), searchQuery)

	if err != nil {
		log.Fatal(err)
	}

	for i, doc := range resDocs {
		fmt.Printf("Doc: %i content: %s", i, doc.PageContent)
	}
}
func textToSplit() []schema.Document {

	f, err := os.Open("./ioc/ioc.go")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}

	p := documentloaders.NewText(f)

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 300   // size of the chunk is number of characters
	split.ChunkOverlap = 30 // overlap is the number of characters that the chunks overlap
	docs, err := p.LoadAndSplit(context.Background(), split)

	if err != nil {
		fmt.Println("Error loading document: ", err)
	}

	log.Println("Document loaded: ", len(docs))

	return docs
}

func ProcessText() {
	ollamaLLM, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatal(err)
	}
	ollamaEmbeder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Qdrant vector store.
	quadrantUrl, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	store, err := qdrant.New(
		qdrant.WithURL(*quadrantUrl),
		qdrant.WithCollectionName(COLLECTION_NAME),
		qdrant.WithEmbedder(ollamaEmbeder),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add documents to the Qdrant vector store.
	_, err = store.AddDocuments(context.Background(), textToSplit())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//ProcessText()
	Search()
}
