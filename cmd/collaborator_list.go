package cmd

import (
	"github.com/spf13/cobra"
)

// lscollabsCmd represents the lscollabs command
var collaboratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collaborators for a repository",
	Long:  `List collaborators for a repository`,
	Run:   ListCollaborators,
}

func init() {
	collaboratorCmd.AddCommand(collaboratorListCmd)
	collaboratorListCmd.PersistentFlags().StringP("repo", "r", "", "specify repository name")
	err := collaboratorListCmd.MarkPersistentFlagRequired("repo")
	if err != nil {
		panic(err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lscollabsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lscollabsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
