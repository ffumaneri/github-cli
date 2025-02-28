package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// aiCmd represents the AI command
var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI requests management.",
	Long: `AI requests management. For example:
git-cli ai generate -p "Generate a README template"
git-cli ai summarize -f myfile.go`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must specify an AI action")
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
