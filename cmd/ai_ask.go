package cmd

import (
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
	Run: AskLlm,
}

func init() {
	aiCmd.AddCommand(askCmd)
}
