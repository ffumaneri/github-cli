package common

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"log"
)

func AskLlm(prompt string) (err error) {
	err = nil
	llm, ok := GetOllamaClient()
	if !ok {
		return nil
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, prompt,
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
