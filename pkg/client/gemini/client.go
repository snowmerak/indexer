package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Config struct {
	APIKey string
	Model  string
}

func NewConfig(apiKey string, model string) *Config {
	return &Config{
		APIKey: apiKey,
		Model:  model,
	}
}

type Client struct {
	client *genai.Client
	config *Config
}

func New(ctx context.Context, opt *Config) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(opt.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	context.AfterFunc(ctx, func() {
		client.Close()
	})

	return &Client{
		client: client,
		config: opt,
	}, nil
}
