/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// lsrepoCmd represents the lsrepo command
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Repository management.",
	Long: `Repository management. For example:
git-cli repository list
git-cli repository invite -r my-repo -c my-collaborator`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must specify a repository action")
	},
}

func init() {
	rootCmd.AddCommand(repositoryCmd)
}
