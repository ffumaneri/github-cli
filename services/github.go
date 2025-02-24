package services

import (
	"fmt"
	github2 "github.com/ffumaneri/github-cli/github"
)

type IGithubService interface {
	ListRepos() error
	ListCollaboratorsByRepo(repo string) error
	InviteCollaboratorToRepo(repo, user string) error
}

func NewGithubService(owner string, githubWrapper github2.IGithubWrapper, consumer func(data string)) *GithubService {
	return &GithubService{
		owner:         owner,
		consumerFunc:  consumer,
		githubWrapper: githubWrapper,
	}
}

type GithubService struct {
	owner         string
	consumerFunc  func(data string)
	githubWrapper github2.IGithubWrapper
}

func (service *GithubService) ListRepos() (err error) {
	repos, err := service.githubWrapper.GetRepos(service.owner)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		service.consumerFunc(repo)
	}
	return
}

func (service *GithubService) ListCollaboratorsByRepo(repo string) (err error) {
	userNames, err := service.githubWrapper.GetCollaboratorsByRepo(service.owner, repo)
	if err != nil {
		return err
	}
	for _, userName := range userNames {
		service.consumerFunc(userName)
	}
	return
}

func (service *GithubService) InviteCollaboratorToRepo(repo, user string) (err error) {
	err = service.githubWrapper.InviteCollaborator(service.owner, repo, user)
	if err != nil {
		return
	}
	service.consumerFunc(fmt.Sprintf("Collaborator %s invited to %s\n", user, repo))
	return
}
