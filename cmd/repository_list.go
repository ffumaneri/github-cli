package cmd

import (
	"github.com/spf13/cobra"
)

// lsrepoCmd represents the lsrepo command
var repositoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories.",
	Long: `List Repositories. For example:
git-cli repository list
`,
	Run: ListRepositories,
}

func init() {
	repositoryCmd.AddCommand(repositoryListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsrepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsrepoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
