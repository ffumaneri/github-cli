package ollama

import (
	"context"
	"errors"
	"testing"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type mockLLM struct {
	CallFunc func(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error)
}

func (m *mockLLM) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	return m.CallFunc(ctx, prompt, opts...)
}

func TestOllamaWrapper_AskLlm(t *testing.T) {
	tests := []struct {
		name          string
		mockCallFunc  func(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error)
		prompt        string
		streamingFunc func(ctx context.Context, chunk []byte) error
		expectError   bool
	}{
		{
			name: "valid prompt with streaming",
			mockCallFunc: func(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
				_ = opts // simulate processing options
				return "response", nil
			},
			prompt: "Hello LLM",
			streamingFunc: func(ctx context.Context, chunk []byte) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "error from LLM",
			mockCallFunc: func(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
				return "", errors.New("llm error")
			},
			prompt: "Hello LLM",
			streamingFunc: func(ctx context.Context, chunk []byte) error {
				return nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLlm := &mockLLM{
				CallFunc: tt.mockCallFunc,
			}
			wrapper := &OllamaWrapper{llm: mockLlm}
			err := wrapper.AskLlm(tt.prompt, tt.streamingFunc)

			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error: got %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestNewOllamaWrapper(t *testing.T) {
	tests := []struct {
		name string
		llm  *ollama.LLM
	}{
		{
			name: "valid LLM",
			llm:  &ollama.LLM{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewOllamaWrapper(tt.llm)
			if wrapper.llm != tt.llm {
				t.Errorf("wrapper.llm = %v, want %v", wrapper.llm, tt.llm)
			}
		})
	}
}
