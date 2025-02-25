package github

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v65/github"
	"github.com/stretchr/testify/assert"
)

type MockGithubRepositories struct {
	mockListByUser        func(ctx context.Context, owner string, opt *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error)
	mockListCollaborators func(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error)
	mockAddCollaborator   func(ctx context.Context, owner, repo, user string, opts *github.RepositoryAddCollaboratorOptions) (*github.CollaboratorInvitation, *github.Response, error)
}

func (m *MockGithubRepositories) ListByUser(ctx context.Context, owner string, opt *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error) {
	return m.mockListByUser(ctx, owner, opt)
}

func (m *MockGithubRepositories) ListCollaborators(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error) {
	return m.mockListCollaborators(ctx, owner, repo, opts)
}

func (m *MockGithubRepositories) AddCollaborator(ctx context.Context, owner, repo, user string, opts *github.RepositoryAddCollaboratorOptions) (*github.CollaboratorInvitation, *github.Response, error) {
	return m.mockAddCollaborator(ctx, owner, repo, user, opts)
}

func TestGetRepos(t *testing.T) {
	tests := []struct {
		name      string
		owner     string
		mockData  []*github.Repository
		mockError error
		want      []string
		wantErr   bool
	}{
		{
			name:  "valid repos list",
			owner: "owner1",
			mockData: []*github.Repository{
				{Name: github.String("repo1"), FullName: github.String("owner1/repo1")},
				{Name: github.String("repo2"), FullName: github.String("owner1/repo2")},
			},
			mockError: nil,
			want:      []string{"owner1/repo1", "owner1/repo2"},
			wantErr:   false,
		},
		{
			name:      "error fetching repos",
			owner:     "owner1",
			mockData:  nil,
			mockError: errors.New("failed to fetch repos"),
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockGithubRepositories{
				mockListByUser: func(ctx context.Context, owner string, opt *github.RepositoryListByUserOptions) ([]*github.Repository, *github.Response, error) {
					return tt.mockData, nil, tt.mockError
				},
			}
			gw := &GithubWrapper{Repositories: mockRepo}
			got, err := gw.GetRepos(tt.owner)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetCollaboratorsByRepo(t *testing.T) {
	tests := []struct {
		name      string
		owner     string
		repo      string
		mockData  []*github.User
		mockError error
		want      []string
		wantErr   bool
	}{
		{
			name:  "valid collaborators list",
			owner: "owner1",
			repo:  "repo1",
			mockData: []*github.User{
				{Login: github.String("user1")},
				{Login: github.String("user2")},
			},
			mockError: nil,
			want:      []string{"user1", "user2"},
			wantErr:   false,
		},
		{
			name:      "error fetching collaborators",
			owner:     "owner1",
			repo:      "repo1",
			mockData:  nil,
			mockError: errors.New("failed to fetch collaborators"),
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockGithubRepositories{
				mockListCollaborators: func(ctx context.Context, owner, repo string, opts *github.ListCollaboratorsOptions) ([]*github.User, *github.Response, error) {
					return tt.mockData, nil, tt.mockError
				},
			}
			gw := &GithubWrapper{Repositories: mockRepo}
			got, err := gw.GetCollaboratorsByRepo(tt.owner, tt.repo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInviteCollaborator(t *testing.T) {
	tests := []struct {
		name      string
		owner     string
		repo      string
		user      string
		mockError error
		wantErr   bool
	}{
		{
			name:      "successful invitation",
			owner:     "owner1",
			repo:      "repo1",
			user:      "user1",
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "failed invitation",
			owner:     "owner1",
			repo:      "repo1",
			user:      "user1",
			mockError: errors.New("failed to invite collaborator"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockGithubRepositories{
				mockAddCollaborator: func(ctx context.Context, owner, repo, user string, opts *github.RepositoryAddCollaboratorOptions) (*github.CollaboratorInvitation, *github.Response, error) {
					return nil, nil, tt.mockError
				},
			}
			gw := &GithubWrapper{Repositories: mockRepo}
			err := gw.InviteCollaborator(tt.owner, tt.repo, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
