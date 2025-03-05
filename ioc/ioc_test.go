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

// NewOllamaService returns a mocked LangChainService.
func (m *MockContainer) NewOllamaService() *services.ILangChainService {
	args := m.Called()
	return args.Get(0).(*services.ILangChainService)
}

func TestMockContainer(t *testing.T) {
	// Create mock objects
	mockGithubService := new(services.IGithubService)    // Real service, customize as necessary
	mockOllamaService := new(services.ILangChainService) // Real service, customize as necessary

	mockContainer := new(MockContainer)

	// Define behavior of the mocked Container
	mockContainer.On("NewGithubService").Return(mockGithubService)
	mockContainer.On("NewOllamaService").Return(mockOllamaService)

	t.Run("GithubServiceCall", func(t *testing.T) {
		// Use the mock container in the test for NewGithubService
		githubService := mockContainer.NewGithubService()
		assert.NotNil(t, githubService)
		mockContainer.AssertCalled(t, "NewGithubService")
	})

	t.Run("LangChainServiceCall", func(t *testing.T) {
		// Use the mock container in the test for NewLangChainService
		ollamaService := mockContainer.NewOllamaService()
		assert.NotNil(t, ollamaService)
		mockContainer.AssertCalled(t, "NewOllamaService")
	})

	t.Run("MockCallTimes", func(t *testing.T) {
		// Ensure methods are called exactly once
		mockContainer.NewGithubService()
		mockContainer.NewOllamaService()
		mockContainer.AssertNumberOfCalls(t, "NewGithubService", 2)
		mockContainer.AssertNumberOfCalls(t, "NewOllamaService", 2)
	})

	t.Run("NoUnexpectedCalls", func(t *testing.T) {
		// Ensure no unexpected methods were called
		mockContainer.AssertExpectations(t)
	})
}
