package text

import (
	"context"
)

type Result[T any] struct {
	Id      int     `json:"id"`
	Payload T       `json:"payload"`
	Score   float64 `json:"score"`
}

type SearchOption struct {
	Limit          int
	Offset         int
	Page           int
	Sort           []string
	ScoreThreshold float64
}

type Text[T any] interface {
	Create(ctx context.Context) error
	Store(ctx context.Context, id int, payload T) error
	Query(ctx context.Context, query string, option SearchOption) ([]Result[T], error)
	Delete(ctx context.Context, id int) error
	Drop(ctx context.Context) error
}
