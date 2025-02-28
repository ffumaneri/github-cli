package cmd

import "github.com/spf13/cobra"

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load source code for analysis",
	Long: `Use this command to load source code into the AI for analysis.
For example:

  load /path/to/sourcefile.go

The AI will analyze the provided source file and process it as needed.`,
	Run: LoadSourceCode,
}

func init() {
	// Define the --name flag and mark it as required
	loadCmd.Flags().StringP("name", "n", "", "Name to associate with the loaded code (required)")
	err := loadCmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	loadCmd.Flags().StringP("path", "p", "", "Path of loading code (required)")
	err = loadCmd.MarkFlagRequired("path")
	if err != nil {
		panic(err)
	}

	aiCmd.AddCommand(loadCmd)
}
