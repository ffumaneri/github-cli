package services

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockOllamaWrapper struct {
	mockAskLlm func(prompt string, handler func(ctx context.Context, chunk []byte) error) error
}

func (m *mockOllamaWrapper) AskLlm(prompt string, handler func(ctx context.Context, chunk []byte) error) error {
	return m.mockAskLlm(prompt, handler)
}

func TestOllamaService_AskLlm(t *testing.T) {
	tests := []struct {
		name          string
		prompt        string
		mockResponse  [][]byte
		mockError     error
		expectedError error
	}{
		{
			name:         "success_single_chunk",
			prompt:       "What is AI?",
			mockResponse: [][]byte{[]byte("AI stands for Artificial Intelligence.")},
		},
		{
			name:         "success_multiple_chunks",
			prompt:       "Explain AI",
			mockResponse: [][]byte{[]byte("AI "), []byte("is "), []byte("Artificial Intelligence.")},
		},
		{
			name:          "error_from_llm_wrapper",
			prompt:        "Invalid input",
			mockResponse:  nil,
			mockError:     errors.New("failed to process prompt"),
			expectedError: errors.New("failed to process prompt"),
		},
		{
			name:   "success_empty_response",
			prompt: "Empty prompt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var results bytes.Buffer
			mockWrapper := &mockOllamaWrapper{
				mockAskLlm: func(prompt string, handler func(ctx context.Context, chunk []byte) error) error {
					if test.mockError != nil {
						return test.mockError
					}
					for _, chunk := range test.mockResponse {
						err := handler(context.Background(), chunk)
						if err != nil {
							return err
						}
					}
					return nil
				},
			}

			service := NewOllamaService(mockWrapper, func(chunk []byte) {
				results.Write(chunk)
			})

			err := service.AskLlm(test.prompt)

			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(bytes.Join(test.mockResponse, nil)), results.String())
			}
		})
	}
}
