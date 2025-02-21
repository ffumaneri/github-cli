package cmd

import (
	"fmt"
	"github.com/ffumaneri/github-cli/ioc"
	"github.com/spf13/cobra"
	// Replace with the correct import path for `AskLlm`
)

// askCmd represents the IA-asking command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask anything to the AI",
	Long: `Use this command to interact with the AI by asking any question. 
For example:

  ask "What is the capital of France?"

The AI will respond with the appropriate answer.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: You must provide a question to ask the AI")
			return
		}

		question := args[0]
		ollamaService := ioc.NewOllamaService()
		err := ollamaService.AskLlm(question)
		if err != nil {
			fmt.Printf("Error while trying to interact with AI: %s\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Flags and configuration options (if needed) can be defined below
	// e.g., askCmd.Flags().StringP("type", "t", "", "Question type")
}
