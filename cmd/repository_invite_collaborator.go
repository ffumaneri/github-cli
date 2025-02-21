package cmd

import (
	"fmt"
	"github.com/ffumaneri/github-cli/ioc"
	"github.com/spf13/cobra"
)

// lsrepoCmd represents the lsrepo command
var repositoryInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			fmt.Println("Too many arguments.")
		} else {
			repo, _ := cmd.Flags().GetString("repo")
			user, _ := cmd.Flags().GetString("collaborator")
			ghService := ioc.NewGithubService()
			ghService.InviteCollaboratorToRepo(repo, user)
		}
	},
}

func init() {
	repositoryCmd.AddCommand(repositoryInviteCmd)
	repositoryInviteCmd.PersistentFlags().StringP("collaborator", "c", "", "specify collaborator name")
	err := repositoryInviteCmd.MarkPersistentFlagRequired("collaborator")
	if err != nil {
		panic(err)
	}
	repositoryInviteCmd.PersistentFlags().StringP("repo", "r", "", "specify repository name")
	err = collaboratorListCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		panic(err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsrepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsrepoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
