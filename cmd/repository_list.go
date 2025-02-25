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
}
