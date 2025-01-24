package ollama

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/snowmerak/indexer/lib/generation"
)

const (
	DefaultURL = "http://localhost:11434"
)

const (
	ModelDefault          = ModelQwen2o5Coder1o5B
	ModelQwen2o5Coder0o5B = "qwen2.5-coder0.5b"
	ModelQwen2o5Coder1o5B = "qwen2.5-coder1.5b"
	ModelQwen2o5Coder3B   = "qwen2.5-coder3b"
	ModelQwen2o5Coder7B   = "qwen2.5-coder7b"
	ModelQwen2o5Coder14B  = "qwen2.5-coder14b"
	ModelQwen2o5Coder32B  = "qwen2.5-coder32b"
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
