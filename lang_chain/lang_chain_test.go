package lang_chain

import (
	"context"
	"errors"
	"github.com/ffumaneri/github-cli/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"os"
	"testing"
)

// MockModel is a mock implementation of the Model interface.
type mockLLM struct {
	// You can define fields for expected behavior or returned data.
	GenerateContentFunc func(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)
	CallFunc            func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

// GenerateContent is the mock implementation of the GenerateContent method.
func (m *mockLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	if m.GenerateContentFunc != nil {
		return m.GenerateContentFunc(ctx, messages, options...)
	}
	return nil, errors.New("GenerateContent method not implemented in mock")
}

// Call is the mock implementation of the Call method.
func (m *mockLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, prompt, options...)
	}
	return "", errors.New("Call method not implemented in mock")
}

type MockTextSplitter struct {
	// Stores the function behavior you want to simulate in the SplitText mock.
	SplitTextFunc func(text string) ([]string, error)
}

// SplitText calls the mocked implementation of SplitTextFunc.
func (m *MockTextSplitter) SplitText(text string) ([]string, error) {
	if m.SplitTextFunc != nil {
		return m.SplitTextFunc(text)
	}
	return nil, errors.New("SplitTextFunc is not implemented")
}

type MockLoader struct {
	mock.Mock
}

func (m *MockLoader) Load(ctx context.Context) ([]schema.Document, error) {
	args := m.Called(ctx)
	return args.Get(0).([]schema.Document), args.Error(1)
}

func (m *MockLoader) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	args := m.Called(ctx, splitter)
	return args.Get(0).([]schema.Document), args.Error(1)
}

type MockVectorStore struct {
	MockAddDocuments     func(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error)
	MockSimilaritySearch func(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error)
}
type MockEmbedder struct {
	EmbedDocumentsFunc func(ctx context.Context, texts []string) ([][]float32, error)
	EmbedQueryFunc     func(ctx context.Context, text string) ([]float32, error)
}

// EmbedDocuments mocks the EmbedDocuments method.
func (m *MockEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	if m.EmbedDocumentsFunc != nil {
		return m.EmbedDocumentsFunc(ctx, texts)
	}
	return nil, errors.New("EmbedDocumentsFunc not implemented")
}

// EmbedQuery mocks the EmbedQuery method.
func (m *MockEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	if m.EmbedQueryFunc != nil {
		return m.EmbedQueryFunc(ctx, text)
	}
	return nil, errors.New("EmbedQueryFunc not implemented")
}

// AddDocuments calls the mock implementation or returns a default response.
func (m *MockVectorStore) AddDocuments(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {
	if m.MockAddDocuments != nil {
		return m.MockAddDocuments(ctx, docs, options...)
	}
	return nil, errors.New("MockAddDocuments not implemented")
}

// SimilaritySearch calls the mock implementation or returns a default response.
func (m *MockVectorStore) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	if m.MockSimilaritySearch != nil {
		return m.MockSimilaritySearch(ctx, query, numDocuments, options...)
	}
	return nil, errors.New("MockSimilaritySearch not implemented")
}

type mockFS struct {
	walkDirFn func(path string, callback common.WalkDirCallback) error
}

func (m *mockFS) WalkDir(path string, callback common.WalkDirCallback) error {
	return m.walkDirFn(path, callback)
}

type MockChain struct {
	mock.Mock
}

// Call mocks the Call method of the Chain interface.
func (m *MockChain) Call(ctx context.Context, inputs map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	args := m.Called(ctx, inputs, options)
	return args.Get(0).(map[string]any), args.Error(1)
}

// GetMemory mocks the GetMemory method of the Chain interface.
func (m *MockChain) GetMemory() schema.Memory {
	args := m.Called()
	return args.Get(0).(schema.Memory)
}

// GetInputKeys mocks the GetInputKeys method of the Chain interface.
func (m *MockChain) GetInputKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// GetOutputKeys mocks the GetOutputKeys method of the Chain interface.
func (m *MockChain) GetOutputKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestLangChainWrapper_AskLlm(t *testing.T) {
	tests := []struct {
		name     string
		mockLLM  *mockLLM
		prompt   string
		expected error
	}{
		{
			name: "success",
			mockLLM: &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "result", nil
				},
			},
			prompt:   "test prompt",
			expected: nil,
		},
		{
			name: "error",
			mockLLM: &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "", errors.New("llm error")
				},
			},
			prompt:   "test prompt",
			expected: errors.New("llm error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewLangChainWrapper(tt.mockLLM, nil, nil, nil, nil, nil, nil)
			err := wrapper.AskLlm(tt.prompt, func(ctx context.Context, chunk []byte) error {
				return nil
			})
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestLangChainWrapper_LoadSourceCode(t *testing.T) {
	tests := []struct {
		name     string
		mockFS   *mockFS
		expected error
	}{
		{
			name: "success",
			mockFS: &mockFS{
				walkDirFn: func(path string, callback common.WalkDirCallback) error {
					callback("path", 1024)
					return nil
				},
			},
			expected: nil,
		},
		{
			name: "error",
			mockFS: &mockFS{
				walkDirFn: func(path string, callback common.WalkDirCallback) error {
					return errors.New("walk error")
				},
			},
			expected: errors.New("walk error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "result", nil
				}}
			wrapper := NewLangChainWrapper(llm, tt.mockFS, func(llm llms.Model) embeddings.Embedder {
				return &MockEmbedder{func(ctx context.Context, texts []string) ([][]float32, error) {
					return [][]float32{}, nil
				}, func(ctx context.Context, text string) ([]float32, error) {
					return []float32{}, nil
				}}
			}, func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore {
				return &MockVectorStore{MockAddDocuments: func(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {
					return []string{}, nil
				}, MockSimilaritySearch: func(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
					return []schema.Document{}, nil
				}}
			}, func() textsplitter.TextSplitter {
				return &MockTextSplitter{}
			}, func(filePath string, size int64) (documentloaders.Loader, *os.File) {
				l := &MockLoader{}
				l.On("LoadAndSplit", mock.Anything, mock.Anything).Return([]schema.Document{}, nil)
				return l, nil
			}, nil)
			err := wrapper.LoadSourceCode("name", "directory")
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}

func TestLangChainWrapper_AskWithContext(t *testing.T) {
	tests := []struct {
		name      string
		mockLLM   *mockLLM
		mockStore func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore
		query     string
		expected  error
	}{
		{
			name: "success",
			mockLLM: &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "result", nil
				},
			},
			mockStore: func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore {
				return nil
			},
			query:    "test query",
			expected: nil,
		},
		{
			name: "error",
			mockLLM: &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "", errors.New("llm error")
				},
			},
			mockStore: func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore {
				return nil
			},
			query:    "test query",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := &mockLLM{
				CallFunc: func(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
					return "result", nil
				}}
			wrapper := NewLangChainWrapper(llm, nil, func(llm llms.Model) embeddings.Embedder {
				return &MockEmbedder{func(ctx context.Context, texts []string) ([][]float32, error) {
					return [][]float32{}, nil
				}, func(ctx context.Context, text string) ([]float32, error) {
					return []float32{}, nil
				}}
			}, func(embedder embeddings.Embedder, collectionName string) vectorstores.VectorStore {
				return &MockVectorStore{MockAddDocuments: func(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {
					return []string{}, nil
				}, MockSimilaritySearch: func(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
					return []schema.Document{}, nil
				}}
			}, func() textsplitter.TextSplitter {
				return &MockTextSplitter{}
			}, func(filePath string, size int64) (documentloaders.Loader, *os.File) {
				l := &MockLoader{}
				l.On("LoadAndSplit", mock.Anything, mock.Anything).Return([]schema.Document{}, nil)
				return l, nil
			}, func(model llms.Model, retriever vectorstores.Retriever) chains.Chain {
				m := &MockChain{}
				m.On("Call", mock.Anything, mock.Anything, mock.Anything).Return(map[string]any{}, nil)
				return m
			})
			err := wrapper.AskWithContext("contextName", tt.query)
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expected.Error())
			}
		})
	}
}
