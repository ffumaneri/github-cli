package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func reportError(format string, args ...any) {
	fmt.Println("FAIL")
	log.Fatalf(format, args)
}

func AskLlm(cmd *cobra.Command, args []string) {
	if len(args) > 2 {
		reportError("Too many arguments. You can only have one which is the repo name")
	}
	contextName, err := cmd.Flags().GetString("name")

	question, err := cmd.Flags().GetString("question")
	if err != nil || question == "" {
		reportError("Question argument is required")
	}
	ollamaService := appContainer.NewOllamaService()
	if contextName == "" {
		err := ollamaService.AskLlm(question)
		if err != nil {
			reportError("Error while trying to interact with AI: %s\n", err)
		}
	} else {
		err := ollamaService.AskLlmWithContext(contextName, question)
		if err != nil {
			reportError("Error while trying to interact with AI: %s\n", err)
		}
	}

}

func ListCollaborators(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		reportError("Too many arguments. You can only have one which is the repo name")
	} else {
		repo, err := cmd.Flags().GetString("repo")
		if err != nil || repo == "" {
			reportError("repo argument %s\n", err)
		}
		ghService := appContainer.NewGithubService()
		err = ghService.ListCollaboratorsByRepo(repo)
		if err != nil {
			reportError("Error while trying to list collaborators: %s\n", err)
		}
	}
}

func InviteCollaborator(cmd *cobra.Command, args []string) {
	if len(args) > 2 {
		reportError("Too many arguments.")
	} else {
		repo, err := cmd.Flags().GetString("repo")
		if err != nil || repo == "" {
			reportError("Repo argument is required")
		}
		user, err := cmd.Flags().GetString("collaborator")
		if err != nil || user == "" {
			reportError("Collaborator argument is required")
		}
		ghService := appContainer.NewGithubService()
		err = ghService.InviteCollaboratorToRepo(repo, user)
		if err != nil {
			reportError("Error while trying to invite collaborator: %s\n", err)
		}
	}
}

func ListRepositories(_ *cobra.Command, _ []string) {
	ghService := appContainer.NewGithubService()
	err := ghService.ListRepos()
	if err != nil {
		reportError("Error while trying to list repositories: %s\n", err)
	}
}

func LoadSourceCode(cmd *cobra.Command, args []string) {
	if len(args) > 2 {
		reportError("Too many arguments.")
	} else {
		name, err := cmd.Flags().GetString("name")
		if err != nil || name == "" {
			reportError("Repo argument is required")
		}
		path, err := cmd.Flags().GetString("path")
		if err != nil || path == "" {
			reportError("Path argument is required")
		}
		oService := appContainer.NewOllamaService()
		err = oService.LoadSourceCode(name, path)
		if err != nil {
			reportError("Error while trying to load documents: %s\n", err)
		}
	}
}
