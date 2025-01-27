package pyembeddings

import (
	"context"

	"github.com/snowmerak/indexer/lib/generation"
)

var _ generation.Embeddings = &Embeddings{}

type Embeddings struct {
	*Client
}

func NewEmbeddings(ctx context.Context, config *Config) (*Embeddings, error) {
	cli, err := NewClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Embeddings{
		Client: cli,
	}, nil
}

func (c *Embeddings) Embed(ctx context.Context, prompt string) ([]float64, error) {
	return c.Client.Embed(prompt)
}

func (c *Embeddings) Size() (uint64, error) {
	s, err := c.Client.Size()

	return uint64(s), err
}
