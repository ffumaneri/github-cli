package services

import (
	"fmt"
	github2 "github.com/ffumaneri/github-cli/github"
	"github.com/google/go-github/v65/github"
	"os"
)

type GithubService struct {
	client *github.Client
	owner  string
}

func NewGithubService(client *github.Client, owner string) *GithubService {
	return &GithubService{
		client: client,
		owner:  owner,
	}
}

func (service *GithubService) ListRepos() {
	repos, err := github2.GetRepos(service.client, service.owner)
	if err != nil {
		fmt.Printf("Error listing repositories %s\n", err)
		os.Exit(1)
	}
	for _, repo := range repos {
		println(repo.GetFullName())
	}
}

func (service *GithubService) ListCollaboratorsByRepo(repo string) {
	users, err := github2.GetCollaboratorsByRepo(service.client, service.owner, repo)
	if err != nil {
		fmt.Printf("Error listing collaborators %s\n", err)
		os.Exit(1)
	}
	for _, user := range users {
		fmt.Printf("%s\n", user.GetLogin())
		//if user.GetLogin() == owner {
		//	continue
		//}
		//invitation, _, err := client.Repositories.AddCollaborator(context.Background(), owner, "2024-tp2-restapi", user.GetLogin(), nil)
		//if err != nil {
		//	println(err.Error())
		//	panic("error enviando invitacion")
		//}
		//println(invitation.Invitee.GetName())
	}
}

func (service *GithubService) InviteCollaboratorToRepo(repo, user string) {
	err := github2.InviteCollaborator(service.client, service.owner, repo, user)
	if err != nil {
		fmt.Printf("Error listing collaborators %s\n", err)
		os.Exit(1)
	}
}
