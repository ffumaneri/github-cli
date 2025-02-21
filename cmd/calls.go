package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func AskLlm(_ *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: You must provide a question to ask the AI")
		return
	}

	question := args[0]
	ollamaService := appContainer.NewOllamaService()
	err := ollamaService.AskLlm(question)
	if err != nil {
		fmt.Printf("Error while trying to interact with AI: %s\n", err)
		return
	}
}
