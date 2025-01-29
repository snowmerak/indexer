package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/pkg/config"
)

const (
	DefaultURL = "http://localhost:11434"
)

const (
	GenerationModelDefault          = GenerationModelQwen2o5Coder1o5B
	GenerationModelQwen2o5Coder0o5B = "qwen2.5-coder:0.5b"
	GenerationModelQwen2o5Coder1o5B = "qwen2.5-coder:1.5b"
	GenerationModelQwen2o5Coder3B   = "qwen2.5-coder:3b"
	GenerationModelQwen2o5Coder7B   = "qwen2.5-coder:7b"
	GenerationModelQwen2o5Coder14B  = "qwen2.5-coder:14b"
	GenerationModelQwen2o5Coder32B  = "qwen2.5-coder:32b"
	GenerationModelDeepseekR11o5B   = "deepseek-r1:1.5b"
	GenerationModelCodeGemma2B      = "codegemma:2b"
)

const (
	EmbeddingModelBgeM3o5B        = "bge-m3"
	EmbeddingModelMxbaiEmbedLarge = "mxbai-embed-large"
)

const (
	EmbeddingModelBgeM3o5BDim        = 1024
	EmbeddingModelMxbaiEmbedLargeDim = 1024
)

type ClientConfig struct {
	url        *url.URL
	httpClient *http.Client
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		httpClient: &http.Client{},
	}
}

func (c *ClientConfig) SetURL(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	c.url = parsed
	return nil
}

func (c *ClientConfig) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

var _ generation.Text = &TextClient{}

func init() {
	generation.RegisterText("ollama", func(ctx context.Context, cc *config.ClientConfig) (generation.Text, error) {
		cfg := NewClientConfig()
		if err := cfg.SetURL(cc.Host[0]); err != nil {
			return nil, fmt.Errorf("failed to set url: %w", err)
		}

		cfg.SetHTTPClient(&http.Client{
			Timeout: 5 * time.Minute,
		})

		return NewTextClient(ctx, cfg, cc.Model)
	})

	generation.RegisterEmbeddings("ollama", func(ctx context.Context, cc *config.ClientConfig) (generation.Embeddings, error) {
		cfg := NewClientConfig()
		if err := cfg.SetURL(cc.Host[0]); err != nil {
			return nil, fmt.Errorf("failed to set url: %w", err)
		}

		cfg.SetHTTPClient(&http.Client{
			Timeout: 5 * time.Minute,
		})

		return NewEmbeddingsClient(ctx, cfg, cc.Model)
	})
}

type TextClient struct {
	client *api.Client
	model  string
}

func NewTextClient(_ context.Context, cfg *ClientConfig, model string) (*TextClient, error) {
	if cfg.url == nil {
		p, err := url.Parse(DefaultURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default url: %w", err)
		}

		cfg.url = p
	}

	client := api.NewClient(cfg.url, cfg.httpClient)

	return &TextClient{
		client: client,
		model:  model,
	}, nil
}

func (c *TextClient) Generate(ctx context.Context, prompt string) (string, error) {
	builder := new(strings.Builder)
	doneCh := make(chan struct{})
	if err := c.client.Generate(ctx, &api.GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
	}, func(response api.GenerateResponse) error {
		builder.WriteString(response.Response)

		if response.Done {
			close(doneCh)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to generate: %w", err)
	}

	for range doneCh {
	}

	return builder.String(), nil
}

var _ generation.Embeddings = &EmbeddingsClient{}

type EmbeddingsClient struct {
	client *api.Client
	model  string
}

func NewEmbeddingsClient(_ context.Context, cfg *ClientConfig, model string) (*EmbeddingsClient, error) {
	if cfg.url == nil {
		p, err := url.Parse(DefaultURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default url: %w", err)
		}

		cfg.url = p
	}

	client := api.NewClient(cfg.url, cfg.httpClient)

	return &EmbeddingsClient{
		client: client,
		model:  model,
	}, nil
}

func (c *EmbeddingsClient) Embed(ctx context.Context, prompt string) ([]float64, error) {
	embeddings, err := c.client.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  c.model,
		Prompt: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embeddings: %w", err)
	}

	return embeddings.Embedding, nil
}

var NotFoundModelErr = fmt.Errorf("model not found")

func (c *EmbeddingsClient) Size() (uint64, error) {
	switch c.model {
	case EmbeddingModelBgeM3o5B:
		return EmbeddingModelBgeM3o5BDim, nil
	case EmbeddingModelMxbaiEmbedLarge:
		return EmbeddingModelMxbaiEmbedLargeDim, nil
	default:
		return 0, NotFoundModelErr
	}
}
