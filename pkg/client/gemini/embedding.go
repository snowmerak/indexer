package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/pkg/config"
)

const (
	ModelTextEmbedding004             = "text-embedding-004"
	ModelTextMultilingualEmbedding002 = "text-multilingual-embedding-002"
	ModelMultiModalEmbedding          = "multi-modal-embedding"
)

const (
	ModelTextEmbedding004Dimension             = 768
	ModelTextMultilingualEmbedding002Dimension = 768
	ModelMultiModalEmbeddingDimension          = 1408
)

var _ generation.Embeddings = &Client{}

func init() {
	generation.RegisterEmbeddings("gemini", func(ctx context.Context, cc *config.ClientConfig) (generation.Embeddings, error) {
		cfg := NewConfig(cc.ApiKey, cc.Model)

		return New(ctx, cfg)
	})
}

func (c *Client) Embed(ctx context.Context, prompt string) ([]float64, error) {
	resp, err := c.client.EmbeddingModel(c.config.Model).EmbedContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to embed content: %w", err)
	}

	if resp.Embedding == nil {
		return nil, fmt.Errorf("no embedding found")
	}

	r := make([]float64, len(resp.Embedding.Values))
	for i, v := range resp.Embedding.Values {
		r[i] = float64(v)
	}

	return r, nil
}

func (c *Client) Size() (uint64, error) {
	switch c.config.Model {
	case ModelTextEmbedding004:
		return ModelTextEmbedding004Dimension, nil
	case ModelTextMultilingualEmbedding002:
		return ModelTextMultilingualEmbedding002Dimension, nil
	case ModelMultiModalEmbedding:
		return ModelMultiModalEmbeddingDimension, nil
	default:
		return 0, fmt.Errorf("unknown model: %s", c.config.Model)
	}
}
