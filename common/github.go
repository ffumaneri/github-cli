package common

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"os"
)

func GetRepos() ([]*github.Repository, error) {
	// Create Git hub client
	client, owner := GetClient()

	repos, _, err := client.Repositories.ListByUser(context.Background(), owner, nil)
	return repos, err
}

func GetCollaboratorsByRepo(repo string) ([]*github.User, error) {
	client, owner := GetClient()
	users, _, err := client.Repositories.ListCollaborators(context.Background(), owner, repo, nil)

	return users, err
}
func InviteCollaborator(repo, user string) error {
	client, owner := GetClient()
	invitation, _, err := client.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Collaborator %s invited to %s\n", invitation.Invitee.GetName(), repo)
	return nil
}
