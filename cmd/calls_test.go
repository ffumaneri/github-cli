package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ffumaneri/github-cli/services"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"os/exec"
	"testing"
)

// MockOllamaService is a mock implementation of the OllamaService
type MockOllamaService struct {
	mock.Mock
}

// AskLlm mocks the AskLlm method of the OllamaService
func (m *MockOllamaService) AskLlm(prompt string) error {
	args := m.Called(prompt)
	return args.Error(0)
}

// MockGithubService is a mock implementation of the GithubServiceInterface
type MockGithubService struct {
	mock.Mock
}

// ListRepos mocks the `ListRepos` method
func (m *MockGithubService) ListRepos() error {
	m.Called()
	return nil
}

// ListCollaboratorsByRepo mocks the `ListCollaboratorsByRepo` method
func (m *MockGithubService) ListCollaboratorsByRepo(repo string) error {
	m.Called(repo)
	return nil
}

// InviteCollaboratorToRepo mocks the `InviteCollaboratorToRepo` method
func (m *MockGithubService) InviteCollaboratorToRepo(repo, user string) error {
	m.Called(repo, user)
	return nil
}

type MockErrorGithubService struct {
	mock.Mock
}

// ListRepos mocks the `ListRepos` method
func (m *MockErrorGithubService) ListRepos() error {
	m.Called()
	return fmt.Errorf("error while listing repos")
}

// ListCollaboratorsByRepo mocks the `ListCollaboratorsByRepo` method
func (m *MockErrorGithubService) ListCollaboratorsByRepo(repo string) error {
	m.Called(repo)
	return fmt.Errorf("error while listing collaborators")
}

// InviteCollaboratorToRepo mocks the `InviteCollaboratorToRepo` method
func (m *MockErrorGithubService) InviteCollaboratorToRepo(repo, user string) error {
	m.Called(repo, user)
	return fmt.Errorf("error while inviting collaborator")
}

type MockContainer struct {
	mock.Mock
	mockGitHubService services.IGithubService
	mockOllamaServie  services.IOllamaService
}

// NewGithubService returns a mocked GithubService.
func (m *MockContainer) NewGithubService() services.IGithubService {
	return m.mockGitHubService
}

// NewOllamaService returns a mocked OllamaService.
func (m *MockContainer) NewOllamaService() services.IOllamaService {
	return m.mockOllamaServie
}

func RunForkTest(_ *testing.T, testName string) (string, string, error) {
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%v", testName))
	cmd.Env = append(os.Environ(), "FORK=1")

	var stdoutB, stderrB bytes.Buffer
	cmd.Stdout = &stdoutB
	cmd.Stderr = &stderrB

	err := cmd.Run()

	return stdoutB.String(), stderrB.String(), err
}

func TestAskLlm_NoArgs(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{} // No arguments provided

	// Capture the output
	output := captureOutput(func() {
		AskLlm(cmd, args)
	})

	assert.Contains(t, output, "Error: You must provide a question to ask the AI", "Expected error message not found.")
}

// TestAskLlm_Success tests the case where AskLlm is successfully called.
func TestAskLlm_Success(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{"What is AI?"} // Valid question

	mockOllamaService := new(MockOllamaService)
	mockOllamaService.On("AskLlm", "What is AI?").Return(nil)
	appContainer = &MockContainer{mockOllamaServie: mockOllamaService}

	// Capture the output
	output := captureOutput(func() {
		AskLlm(cmd, args)
	})

	assert.Empty(t, output, nil)
	mockOllamaService.AssertCalled(t, "AskLlm", "What is AI?")
}

// TestAskLlm_Error tests the case where AskLlm returns an error.
func TestAskLlm_Error(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		args := []string{"What is AI?"} // Valid question

		mockOllamaService := new(MockOllamaService)
		mockOllamaService.On("AskLlm", "What is AI?").Return(errors.New("failed to query AI"))
		appContainer = &MockContainer{mockOllamaServie: mockOllamaService}

		AskLlm(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestAskLlm_Error")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Error while trying to interact with AI")
	assert.Contains(t, stdout, "FAIL")
}

func TestListCollaborators_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("repo", "test-repo", "Name of the repository") // Set the repo flag
	args := []string{}                                                // No additional arguments passed

	mockGithubService := new(MockGithubService)
	mockGithubService.On("ListCollaboratorsByRepo", "test-repo").Return(nil)

	// Inject the mock service into the app container
	appContainer = &MockContainer{mockGitHubService: mockGithubService}

	// Capture the output
	output := captureOutput(func() {
		ListCollaborators(cmd, args)
	})

	// Since the output in this use case doesn't print success, ensure it's empty
	assert.Empty(t, output, "Expected no output on success")
	mockGithubService.AssertCalled(t, "ListCollaboratorsByRepo", "test-repo")
}

func TestListCollaborators_TooManyArguments(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		args := []string{"arg1", "arg2"} // Too many arguments
		ListCollaborators(cmd, args)
	}
	stdout, stderr, err := RunForkTest(t, "TestListCollaborators_TooManyArguments")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Too many arguments. You can only have one which is the repo name")
	assert.Contains(t, stdout, "FAIL")
}

func TestListCollaborators_RepoFlagEmpty(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		cmd.Flags().String("repo", "", "Name of the repository") // Empty repo flag
		args := []string{}                                       // No additional arguments passed

		// Mock service never gets invoked
		mockGithubService := new(MockGithubService)
		appContainer = &MockContainer{mockGitHubService: mockGithubService}
		ListCollaborators(cmd, args)
	}
	stdout, stderr, err := RunForkTest(t, "TestListCollaborators_RepoFlagEmpty")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "repo argument")
	assert.Contains(t, stdout, "FAIL")
}

func TestListCollaborators_WithError(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		cmd.Flags().String("repo", "test-repo", "Name of the repository") // Set the repo flag
		args := []string{}                                                // No additional arguments passed

		// Mock service never gets invoked
		mockGithubService := new(MockErrorGithubService)
		mockGithubService.On("ListCollaboratorsByRepo", "test-repo").Return(errors.New("error while listing collaborators"))
		appContainer = &MockContainer{mockGitHubService: mockGithubService}
		ListCollaborators(cmd, args)
	}
	stdout, stderr, err := RunForkTest(t, "TestListCollaborators_WithError")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Error while trying to list collaborators")
	assert.Contains(t, stdout, "FAIL")
}
func TestInviteCollaborator_Success(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("repo", "test-repo", "Repository name")
	cmd.Flags().String("collaborator", "test-user", "Collaborator username")
	args := []string{} // No additional arguments passed

	mockGithubService := new(MockGithubService)
	mockGithubService.On("InviteCollaboratorToRepo", "test-repo", "test-user").Return(nil)

	// Inject the mock service into the app container
	appContainer = &MockContainer{mockGitHubService: mockGithubService}

	// Capture the output
	output := captureOutput(func() {
		InviteCollaborator(cmd, args)
	})

	// Since the function doesn't print anything on success, output should be empty
	assert.Empty(t, output, "Expected no output on success")
	mockGithubService.AssertCalled(t, "InviteCollaboratorToRepo", "test-repo", "test-user")
}

func TestInviteCollaborator_WithError(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		cmd.Flags().String("repo", "test-repo", "Repository name")
		cmd.Flags().String("collaborator", "test-user", "Collaborator username")
		args := []string{} // No additional arguments passed

		mockGithubService := new(MockErrorGithubService)
		mockGithubService.On("InviteCollaboratorToRepo", "test-repo", "test-user").Return(fmt.Errorf("error"))

		// Inject the mock service into the app container
		appContainer = &MockContainer{mockGitHubService: mockGithubService}

		InviteCollaborator(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestInviteCollaborator_WithError")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Error while trying to invite collaborator")
	assert.Contains(t, stdout, "FAIL")
}
func TestInviteCollaborator_TooManyArguments(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		args := []string{"arg1", "arg2", "arg3"} // Too many arguments

		// Capture the output
		InviteCollaborator(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestInviteCollaborator_TooManyArguments")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Too many arguments")
	assert.Contains(t, stdout, "FAIL")

}

func TestInviteCollaborator_MissingRepoFlag(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		cmd.Flags().String("repo", "", "Repository name") // Empty repo flag
		cmd.Flags().String("collaborator", "test-user", "Collaborator username")
		args := []string{} // No additional arguments

		InviteCollaborator(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestInviteCollaborator_MissingRepoFlag")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Repo argument is required")
	assert.Contains(t, stdout, "FAIL")
}

func TestInviteCollaborator_MissingCollaboratorFlag(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		cmd.Flags().String("repo", "test-repo", "Repository name")
		cmd.Flags().String("collaborator", "", "Collaborator username") // Empty collaborator flag
		args := []string{}                                              // No additional arguments

		InviteCollaborator(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestInviteCollaborator_MissingCollaboratorFlag")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Collaborator argument is required")
	assert.Contains(t, stdout, "FAIL")
}

func TestListRepositories_Success(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{} // No arguments expected

	mockGithubService := new(MockGithubService)
	mockGithubService.On("ListRepos").Return(nil)

	// Inject the mock service into the app container
	appContainer = &MockContainer{mockGitHubService: mockGithubService}

	// Capture the output
	output := captureOutput(func() {
		ListRepositories(cmd, args)
	})

	// Since the function doesn't print anything on success, output should be empty
	assert.Empty(t, output, "Expected no output on success")
	mockGithubService.AssertCalled(t, "ListRepos")
}

func TestListRepositories_WithError(t *testing.T) {
	if os.Getenv("FORK") == "1" {
		cmd := &cobra.Command{}
		args := []string{} // No arguments expected

		mockGithubService := new(MockErrorGithubService)
		mockGithubService.On("ListRepos").Return(fmt.Errorf("error"))

		// Inject the mock service into the app container
		appContainer = &MockContainer{mockGitHubService: mockGithubService}

		ListRepositories(cmd, args)
	}

	stdout, stderr, err := RunForkTest(t, "TestListRepositories_WithError")

	assert.NotNil(t, err, "Expected error not found.")
	assert.Equal(t, err.Error(), "exit status 1")
	assert.Contains(t, stderr, "Error while trying to list repositories")
	assert.Contains(t, stdout, "FAIL")
}

// Helper function to capture console output for testing
func captureOutput(f func()) string {
	// Create a pipe to redirect os.Stdout
	r, w, _ := os.Pipe()
	// Save the original os.Stdout for restoration later
	stdout := os.Stdout
	// Set os.Stdout to the write end of the pipe
	os.Stdout = w

	// Call the function whose output we want to capture
	f()

	// Close the writer and restore os.Stdout
	_ = w.Close()
	os.Stdout = stdout

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	_ = r.Close()

	return buf.String()
}
