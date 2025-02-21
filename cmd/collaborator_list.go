package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// lscollabsCmd represents the lscollabs command
var collaboratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collaborators for a repository",
	Long:  `List collaborators for a repository`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments. You can only have one which is the repo name")
			os.Exit(1)
		} else {
			repo, err := cmd.Flags().GetString("repo")
			if err != nil {
				fmt.Printf("repo argument %s\n", err)
				os.Exit(1)
			}
			ghService := AppContainer.NewGithubService()
			ghService.ListCollaboratorsByRepo(repo)
		}
	},
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
