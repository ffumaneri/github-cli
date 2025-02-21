package cmd

import (
	"github.com/ffumaneri/github-cli/ioc"
	"github.com/spf13/cobra"
)

// lsrepoCmd represents the lsrepo command
var repositoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ghService := ioc.NewGithubService()
		ghService.ListRepos()
	},
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
