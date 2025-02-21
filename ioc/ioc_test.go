package ioc

import (
	"github.com/ffumaneri/github-cli/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockContainer is a mocked implementation of the Container interface using stretchr/testify/mock.
type MockContainer struct {
	mock.Mock
}

// NewGithubService returns a mocked GithubService.
func (m *MockContainer) NewGithubService() *services.IGithubService {
	args := m.Called()
	return args.Get(0).(*services.IGithubService)
}

// NewOllamaService returns a mocked OllamaService.
func (m *MockContainer) NewOllamaService() *services.IOllamaService {
	args := m.Called()
	return args.Get(0).(*services.IOllamaService)
}

func TestMockContainer(t *testing.T) {
	// Create mock objects
	mockGithubService := new(services.IGithubService) // Real service, customize as necessary
	mockOllamaService := new(services.IOllamaService) // Real service, customize as necessary

	mockContainer := new(MockContainer)

	// Define behavior of the mocked Container
	mockContainer.On("NewGithubService").Return(mockGithubService)
	mockContainer.On("NewOllamaService").Return(mockOllamaService)

	// Use the mock container in the tests
	githubService := mockContainer.NewGithubService()
	ollamaService := mockContainer.NewOllamaService()

	// Assertions
	assert.NotNil(t, githubService)
	assert.NotNil(t, ollamaService)

	// Verify that the expected methods were called
	mockContainer.AssertCalled(t, "NewGithubService")
	mockContainer.AssertCalled(t, "NewOllamaService")
}
