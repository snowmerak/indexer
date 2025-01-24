package generation

import "context"

type Text interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
