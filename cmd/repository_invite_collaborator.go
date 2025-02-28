package cmd

import (
	"github.com/spf13/cobra"
)

// lsrepoCmd represents the lsrepo command
var repositoryInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a collaborator to a repository.",
	Long: `Invite a collaborator to a repository. For example:
git-cli repository invite -r my-repo -c my-collaborator
`,
	Run: InviteCollaborator,
}

func init() {
	repositoryCmd.AddCommand(repositoryInviteCmd)
	repositoryInviteCmd.PersistentFlags().StringP("collaborator", "c", "", "specify collaborator name")
	err := repositoryInviteCmd.MarkPersistentFlagRequired("collaborator")
	if err != nil {
		panic(err)
	}
	repositoryInviteCmd.PersistentFlags().StringP("repo", "r", "", "specify repository name")
	err = CollaboratorListCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		panic(err)
	}
}
