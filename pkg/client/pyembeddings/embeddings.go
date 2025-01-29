package pyembeddings

import (
	"context"
	"net/http"
	"time"

	"github.com/snowmerak/indexer/lib/generation"
	"github.com/snowmerak/indexer/pkg/config"
)

var _ generation.Embeddings = &Embeddings{}

func init() {
	generation.RegisterEmbeddings("pyembeddings", func(ctx context.Context, cc *config.ClientConfig) (generation.Embeddings, error) {
		cfg := NewConfig(cc.Host[0]).
			WithHttpClient(&http.Client{
				Timeout: 1 * time.Minute,
			})

		return NewEmbeddings(ctx, cfg)
	})
}

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
