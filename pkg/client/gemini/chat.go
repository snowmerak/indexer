package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/pkg/config"
)

const (
	ModelGemini2o0Flash   = "gemini-2.0-flash"
	ModelGemini1o5Flash   = "gemini-1.5-flash"
	ModelGemini1o5Flash8B = "gemini-1.5-flash-8b"
	ModelGemini1o5Pro     = "gemini-1.5-pro"
	ModelAQA              = "aqa"
)

var _ generation.Text = &Client{}

func init() {
	generation.RegisterText("gemini", func(ctx context.Context, cc *config.ClientConfig) (generation.Text, error) {
		cfg := NewConfig(cc.ApiKey, cc.Model)

		return New(ctx, cfg)
	})
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.GenerativeModel(c.config.Model).GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	builder := &strings.Builder{}
	for _, c := range resp.Candidates {
		if c == nil {
			continue
		}

		if c.Content == nil {
			continue
		}

		for _, p := range c.Content.Parts {
			t, ok := p.(*genai.Text)
			if !ok {
				continue
			}

			builder.WriteString(string(*t))
		}
	}

	return builder.String(), nil
}
