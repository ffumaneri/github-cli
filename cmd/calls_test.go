package cmd

import (
	"bytes"
	"errors"
	"github.com/ffumaneri/github-cli/services"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
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
func (m *MockGithubService) ListRepos() {
	m.Called()
}

// ListCollaboratorsByRepo mocks the `ListCollaboratorsByRepo` method
func (m *MockGithubService) ListCollaboratorsByRepo(repo string) {
	m.Called(repo)
}

// InviteCollaboratorToRepo mocks the `InviteCollaboratorToRepo` method
func (m *MockGithubService) InviteCollaboratorToRepo(repo, user string) {
	m.Called(repo, user)
}

type MockContainer struct {
	mock.Mock
	mockGitHubService *MockGithubService
	mockOllamaServie  *MockOllamaService
}

// NewGithubService returns a mocked GithubService.
func (m *MockContainer) NewGithubService() services.IGithubService {
	return new(MockGithubService)
}

// NewOllamaService returns a mocked OllamaService.
func (m *MockContainer) NewOllamaService() services.IOllamaService {
	return m.mockOllamaServie
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
	cmd := &cobra.Command{}
	args := []string{"What is AI?"} // Valid question

	mockOllamaService := new(MockOllamaService)
	mockOllamaService.On("AskLlm", "What is AI?").Return(errors.New("failed to query AI"))
	appContainer = &MockContainer{mockOllamaServie: mockOllamaService}

	// Capture the output
	output := captureOutput(func() {
		AskLlm(cmd, args)
	})

	// Verify the error message is printed
	assert.Contains(t, output, "Error while trying to interact with AI: failed to query AI", "Expected error message not found.")
	mockOllamaService.AssertCalled(t, "AskLlm", "What is AI?")
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
