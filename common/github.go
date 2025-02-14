package common

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"os"
)

func GetRepos(ctx context.Context) ([]*github.Repository, error) {
	client := ctx.Value("githubClient").(*github.Client)
	owner := ctx.Value("githubOwner").(string)

	repos, _, err := client.Repositories.ListByUser(context.Background(), owner, nil)
	return repos, err
}

func GetCollaboratorsByRepo(ctx context.Context, repo string) ([]*github.User, error) {
	client := ctx.Value("githubClient").(*github.Client)
	owner := ctx.Value("githubOwner").(string)

	users, _, err := client.Repositories.ListCollaborators(context.Background(), owner, repo, nil)
	return users, err
}
func InviteCollaborator(ctx context.Context, repo, user string) error {
	client := ctx.Value("githubClient").(*github.Client)
	owner := ctx.Value("githubOwner").(string)

	invitation, _, err := client.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Collaborator %s invited to %s\n", invitation.Invitee.GetName(), repo)
	return nil
}
