package text

import (
	"context"
)

type Payload struct {
	Description string `json:"description"`
	CodeBlock   string `json:"code_block"`
}

type Result struct {
	Id      int     `json:"id"`
	Payload Payload `json:"payload"`
	Score   float64 `json:"score"`
}

type SearchOption struct {
	Limit          int
	Offset         int
	Page           int
	Sort           []string
	ScoreThreshold float64
}

type Text interface {
	Create(ctx context.Context) error
	Store(ctx context.Context, id int, payload Payload) error
	Query(ctx context.Context, query string, option SearchOption) ([]Result, error)
	Delete(ctx context.Context, id int) error
	Drop(ctx context.Context) error
	UpdateSynonyms(ctx context.Context, synonyms map[string][]string) error
}
