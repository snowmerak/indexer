package vector

import "context"

type Vector interface {
	NextId() int
	Store(ctx context.Context, id int, vector []float64) error
	Get(ctx context.Context, id int) ([]float64, error)
	Search(ctx context.Context, vector []float64, limit int) ([]int, error)
	Delete(ctx context.Context, id int) error
	Drop(ctx context.Context) error
}
