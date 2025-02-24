package services

import (
	"fmt"
	github2 "github.com/ffumaneri/github-cli/github"
	"github.com/google/go-github/v65/github"
)

type IGithubService interface {
	ListRepos() error
	ListCollaboratorsByRepo(repo string) error
	InviteCollaboratorToRepo(repo, user string) error
}

type GithubService struct {
	client       *github.Client
	owner        string
	consumerFunc func(data string)
}

func NewGithubService(client *github.Client, owner string, consumer func(data string)) *GithubService {
	return &GithubService{
		client:       client,
		owner:        owner,
		consumerFunc: consumer,
	}
}

func (service *GithubService) ListRepos() (err error) {
	repos, err := github2.GetRepos(service.client, service.owner)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		service.consumerFunc(repo)
	}
	return
}

func (service *GithubService) ListCollaboratorsByRepo(repo string) (err error) {
	userNames, err := github2.GetCollaboratorsByRepo(service.client, service.owner, repo)
	if err != nil {
		return err
	}
	for _, userName := range userNames {
		service.consumerFunc(userName)
	}
	return
}

func (service *GithubService) InviteCollaboratorToRepo(repo, user string) (err error) {
	err = github2.InviteCollaborator(service.client, service.owner, repo, user)
	if err != nil {
		return
	}
	service.consumerFunc(fmt.Sprintf("Collaborator %s invited to %s\n", user, repo))
	return
}
