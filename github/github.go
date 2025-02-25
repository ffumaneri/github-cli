package github

import (
	"context"
	"github.com/google/go-github/v65/github"
)

type IGithubWrapper interface {
	GetRepos(owner string) ([]string, error)
	GetCollaboratorsByRepo(owner string, repo string) ([]string, error)
	InviteCollaborator(owner string, repo, user string) error
}

func NewGithubWrapper(client *github.Client, owner string) *GithubWrapper {
	return &GithubWrapper{client.Repositories, owner}
}

type IGithubRepositories interface {
	ListByUser(ctx context.Context, owner string, opt *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error)
	AddCollaborator(ctx context.Context, owner, repo, user string, opts *github.RepositoryAddCollaboratorOptions) (*github.CollaboratorInvitation, *github.Response, error)
	ListCollaborators(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error)
}
type GithubWrapper struct {
	Repositories IGithubRepositories
	owner        string
}

func (gw *GithubWrapper) GetRepos(owner string) ([]string, error) {
	repos, _, err := gw.Repositories.ListByUser(context.Background(), owner, nil)
	repoNames := make([]string, len(repos))
	for i, repo := range repos {
		repoNames[i] = repo.GetFullName()
	}
	return repoNames, err
}

func (gw *GithubWrapper) GetCollaboratorsByRepo(owner string, repo string) ([]string, error) {
	users, _, err := gw.Repositories.ListCollaborators(context.Background(), owner, repo, nil)
	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.GetLogin()
	}
	return userNames, err
}
func (gw *GithubWrapper) InviteCollaborator(owner string, repo, user string) error {
	_, _, err := gw.Repositories.AddCollaborator(context.Background(), owner, repo, user, nil)
	if err != nil {
		return err
	}
	return nil
}
