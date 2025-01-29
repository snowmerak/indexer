package text

import (
	"context"
	"fmt"
	"sync"

	"github.com/snowmerak/indexer/pkg/config"
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

var registeredText = sync.Map{}

type TextConstructor func(context.Context, *config.ClientConfig) (Text, error)

func RegisterText(name string, text TextConstructor) {
	registeredText.Store(name, text)
}

func GetText(ctx context.Context, name string, config *config.ClientConfig) (Text, error) {
	if v, ok := registeredText.Load(name); ok {
		text, err := v.(TextConstructor)(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("failed to create text client: %w", err)
		}
		return text, nil
	}
	return nil, fmt.Errorf("text client not found: %s", name)
}
