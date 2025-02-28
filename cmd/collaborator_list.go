package cmd

import (
	"github.com/spf13/cobra"
)

// lscollabsCmd represents the lscollabs command
var CollaboratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collaborators for a repository",
	Long:  `List collaborators for a repository`,
	Run:   ListCollaborators,
}

func init() {
	collaboratorCmd.AddCommand(CollaboratorListCmd)
	CollaboratorListCmd.PersistentFlags().StringP("repo", "r", "", "specify repository name")
	err := CollaboratorListCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		panic(err)
	}
}
