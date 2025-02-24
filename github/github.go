package github

import (
	"context"
	"github.com/ffumaneri/github-cli/common"
	"github.com/google/go-github/v65/github"
)

type IGithubWrapper interface {
	GetRepos(owner string) ([]string, error)
	GetCollaboratorsByRepo(owner string, repo string) ([]string, error)
	InviteCollaborator(owner string, repo, user string) error
}

func NewGithubClient(config *common.Config) (*github.Client, string, error) {
	// Create Github client
	client := github.NewClient(nil).WithAuthToken(config.Token)
	return client, config.Owner, nil
}
func NewGithubWrapper(client *github.Client, owner string) *GithubWrapper {
	return &GithubWrapper{client, owner}
}

type GithubWrapper struct {
	client *github.Client
	owner  string
}

func (gw *GithubWrapper) GetRepos(owner string) ([]string, error) {
	repos, _, err := gw.client.Repositories.ListByUser(context.Background(), owner, nil)
	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = repo.GetFullName()
	}
	return repoNames, err
}

func (gw *GithubWrapper) GetCollaboratorsByRepo(owner string, repo string) ([]string, error) {
	users, _, err := gw.client.Repositories.ListCollaborators(context.Background(), owner, repo, nil)
	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.GetLogin()
	}
	return userNames, err
}
func (gw *GithubWrapper) InviteCollaborator(owner string, repo, user string) error {
	_, _, err := gw.client.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		return err
	}
	return nil
}
