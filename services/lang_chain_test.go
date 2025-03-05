package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockLangChainWrapper struct {
	mock.Mock
}

func (m *MockLangChainWrapper) AskLlm(prompt string, streamingFunc func(ctx context.Context, chunk []byte) error) error {
	args := m.Called(prompt, streamingFunc)
	if streamingFunc != nil {
		// If a streamingFunc is provided, invoke it with a sample context and chunk for testing.
		streamingFunc(context.TODO(), []byte("mock chunk data"))
	}
	return args.Error(0)
}

func (m *MockLangChainWrapper) LoadSourceCode(name, directory string) error {
	args := m.Called(name, directory)
	return args.Error(0)
}

func (m *MockLangChainWrapper) AskWithContext(contextName, searchQuery string) error {
	args := m.Called(contextName, searchQuery)
	return args.Error(0)
}
func TestLangChainService_AskLlm(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		mockError      error
		expectedError  error
		streamingError error
		consumerChunks []string
	}{
		{
			name:           "success_valid_prompt",
			prompt:         "What is AI?",
			consumerChunks: []string{"mock chunk data", "mock chunk data"},
		},
		{
			name:          "error_empty_prompt",
			prompt:        "",
			mockError:     errors.New("prompt cannot be empty"),
			expectedError: errors.New("prompt cannot be empty"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockWrapper := &MockLangChainWrapper{}
			mockWrapper.On("AskLlm", test.prompt, mock.Anything).Return(test.mockError).Run(func(args mock.Arguments) {
				streamingFunc := args.Get(1).(func(ctx context.Context, chunk []byte) error)
				if streamingFunc != nil {
					if test.streamingError != nil {
						assert.EqualError(t, streamingFunc(context.TODO(), []byte("mock chunk data")), test.streamingError.Error())
					} else {
						assert.NoError(t, streamingFunc(context.TODO(), []byte("mock chunk data")))
					}
				}
			})

			var consumerChunks []string
			service := NewLangChainService(mockWrapper, func(chunk []byte) {
				consumerChunks = append(consumerChunks, string(chunk))
			})

			err := service.AskLlm(test.prompt)

			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.consumerChunks, consumerChunks)
			}
		})
	}
}
func TestLangChainService_AskLlmWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextName   string
		prompt        string
		mockError     error
		expectedError error
	}{
		{
			name:        "success_with_valid_context",
			contextName: "validContext",
			prompt:      "What is the purpose of AI?",
		},
		{
			name:          "error_invalid_context",
			contextName:   "invalidContext",
			prompt:        "Explain ML",
			mockError:     errors.New("invalid context"),
			expectedError: errors.New("invalid context"),
		},
		{
			name:          "error_empty_context",
			contextName:   "",
			prompt:        "Provide details on deep learning",
			mockError:     errors.New("context cannot be empty"),
			expectedError: errors.New("context cannot be empty"),
		},
		{
			name:          "error_empty_prompt",
			contextName:   "someContext",
			prompt:        "",
			mockError:     errors.New("prompt cannot be empty"),
			expectedError: errors.New("prompt cannot be empty"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockWrapper := &MockLangChainWrapper{}
			mockWrapper.On("AskWithContext", test.contextName, test.prompt).Return(test.mockError)

			service := NewLangChainService(mockWrapper, nil)

			err := service.AskLlmWithContext(test.contextName, test.prompt)

			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLangChainService_LoadSourceCode(t *testing.T) {
	tests := []struct {
		name          string
		sourceName    string
		directory     string
		mockError     error
		expectedError error
	}{
		{
			name:       "success_valid_directory",
			sourceName: "exampleSource",
			directory:  "/valid/path",
		},
		{
			name:          "error_invalid_directory",
			sourceName:    "brokenSource",
			directory:     "/invalid/path",
			mockError:     errors.New("directory not found"),
			expectedError: errors.New("directory not found"),
		},
		{
			name:          "error_large_file",
			sourceName:    "largeFile",
			directory:     "/large/file/path",
			mockError:     errors.New("file too large to process"),
			expectedError: errors.New("file too large to process"),
		},
		{
			name:          "error_document_processing",
			sourceName:    "errorDocument",
			directory:     "/path/to/error",
			mockError:     errors.New("failed to process document"),
			expectedError: errors.New("failed to process document"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockWrapper := &MockLangChainWrapper{}
			mockWrapper.On("LoadSourceCode", test.sourceName, test.directory).Return(test.mockError)

			service := NewLangChainService(mockWrapper, nil)

			err := service.LoadSourceCode(test.sourceName, test.directory)

			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
