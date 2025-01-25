package generation

import "context"

type Text interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Embeddings interface {
	Embed(ctx context.Context, prompt string) ([]float64, error)
	Size() uint64
}
