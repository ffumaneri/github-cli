package common

import (
	"fmt"
	"os"
)

func ListRepos() {
	repos, err := GetRepos()
	if err != nil {
		fmt.Printf("Error listing repositories %s\n", err)
		os.Exit(1)
	}
	for _, repo := range repos {
		println(repo.GetFullName())
	}
}

func ListCollaboratorsByRepo(repo string) {
	users, err := GetCollaboratorsByRepo(repo)
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

func InviteCollaboratorToRepo(repo, user string) {
	err := InviteCollaborator(repo, user)
	if err != nil {
		fmt.Printf("Error listing collaborators %s\n", err)
		os.Exit(1)
	}
}
