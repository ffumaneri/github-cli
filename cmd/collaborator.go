/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// lscollabsCmd represents the lscollabs command
var collaboratorCmd = &cobra.Command{
	Use:   "collaborator",
	Short: "Collaborators for a repository.",
	Long:  `This command allows you to manage collaborators for a repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must specify a collaborator action")
	},
}

func init() {
	rootCmd.AddCommand(collaboratorCmd)
}
