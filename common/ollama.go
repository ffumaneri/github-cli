package common

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

func AskLlm(ctx context.Context, prompt string) (err error) {
	llm := ctx.Value("ollamaClient").(*ollama.LLM)

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
