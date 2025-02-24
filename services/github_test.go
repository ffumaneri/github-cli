package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGithubWrapper struct {
	mock.Mock
}

func (m *MockGithubWrapper) GetRepos(owner string) ([]string, error) {
	args := m.Called(owner)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGithubWrapper) GetCollaboratorsByRepo(owner, repo string) ([]string, error) {
	args := m.Called(owner, repo)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGithubWrapper) InviteCollaborator(owner, repo, user string) error {
	args := m.Called(owner, repo, user)
	return args.Error(0)
}

func TestGithubService_ListRepos(t *testing.T) {
	tests := []struct {
		name          string
		mockRepos     []string
		mockError     error
		expectedError error
	}{
		{"Success", []string{"repo1", "repo2"}, nil, nil},
		{"API error", nil, errors.New("API error"), errors.New("API error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWrapper := new(MockGithubWrapper)
			consumerOutput := []string{}
			consumerFunc := func(data string) { consumerOutput = append(consumerOutput, data) }
			service := NewGithubService("owner", mockWrapper, consumerFunc)

			mockWrapper.On("GetRepos", "owner").Return(tt.mockRepos, tt.mockError)

			err := service.ListRepos()

			assert.Equal(t, tt.expectedError, err)
			if err == nil {
				assert.Equal(t, tt.mockRepos, consumerOutput)
			}
			mockWrapper.AssertExpectations(t)
		})
	}
}

func TestGithubService_ListCollaboratorsByRepo(t *testing.T) {
	tests := []struct {
		name          string
		repo          string
		mockUsers     []string
		mockError     error
		expectedError error
	}{
		{"Success", "repo1", []string{"user1", "user2"}, nil, nil},
		{"API error", "repo1", nil, errors.New("API error"), errors.New("API error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWrapper := new(MockGithubWrapper)
			consumerOutput := []string{}
			consumerFunc := func(data string) { consumerOutput = append(consumerOutput, data) }
			service := NewGithubService("owner", mockWrapper, consumerFunc)

			mockWrapper.On("GetCollaboratorsByRepo", "owner", tt.repo).Return(tt.mockUsers, tt.mockError)

			err := service.ListCollaboratorsByRepo(tt.repo)

			assert.Equal(t, tt.expectedError, err)
			if err == nil {
				assert.Equal(t, tt.mockUsers, consumerOutput)
			}
			mockWrapper.AssertExpectations(t)
		})
	}
}

func TestGithubService_InviteCollaboratorToRepo(t *testing.T) {
	tests := []struct {
		name           string
		repo           string
		user           string
		mockError      error
		expectedOutput string
		expectedError  error
	}{
		{"Success", "repo1", "user1", nil, "Collaborator user1 invited to repo1\n", nil},
		{"API error", "repo1", "user1", errors.New("API error"), "", errors.New("API error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWrapper := new(MockGithubWrapper)
			consumerOutput := []string{}
			consumerFunc := func(data string) { consumerOutput = append(consumerOutput, data) }
			service := NewGithubService("owner", mockWrapper, consumerFunc)

			mockWrapper.On("InviteCollaborator", "owner", tt.repo, tt.user).Return(tt.mockError)

			err := service.InviteCollaboratorToRepo(tt.repo, tt.user)

			assert.Equal(t, tt.expectedError, err)
			if err == nil {
				assert.Equal(t, []string{tt.expectedOutput}, consumerOutput)
			}
			mockWrapper.AssertExpectations(t)
		})
	}
}
