package github

import (
	"context"
	"github.com/ffumaneri/github-cli/common"
	"github.com/google/go-github/v65/github"
)

func NewGithubClient(config *common.Config) (*github.Client, string, error) {
	// Create Github client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}

func GetRepos(client *github.Client, owner string) ([]string, error) {
	repos, _, err := client.Repositories.ListByUser(context.Background(), owner, nil)
	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = repo.GetFullName()
	}
	return repoNames, err
}

func GetCollaboratorsByRepo(client *github.Client, owner string, repo string) ([]string, error) {
	users, _, err := client.Repositories.ListCollaborators(context.Background(), owner, repo, nil)
	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.GetLogin()
	}
	return userNames, err
}
func InviteCollaborator(client *github.Client, owner string, repo, user string) error {
	_, _, err := client.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		return err
	}
	return nil
}
