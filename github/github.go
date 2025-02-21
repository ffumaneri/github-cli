package github

import (
	"context"
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"github.com/google/go-github/v65/github"
	"os"
)

func NewGithubClient(config *common.Config) (*github.Client, string, error) {
	// Create Github client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}

func GetRepos(client *github.Client, owner string) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.ListByUser(context.Background(), owner, nil)
	return repos, err
}

func GetCollaboratorsByRepo(client *github.Client, owner string, repo string) ([]*github.User, error) {
	users, _, err := client.Repositories.ListCollaborators(context.Background(), owner, repo, nil)
	return users, err
}
func InviteCollaborator(client *github.Client, owner string, repo, user string) error {
	invitation, _, err := client.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Collaborator %s invited to %s\n", invitation.Invitee.GetName(), repo)
	return nil
}
