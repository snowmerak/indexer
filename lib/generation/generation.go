package generation

import (
	"context"
	"fmt"
	"sync"

	"github.com/snowmerak/indexer/pkg/config"
)

type Text interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

var registeredText = sync.Map{}

type TextConstructor func(*config.ClientConfig) (Text, error)

func RegisterText(name string, text TextConstructor) {
	registeredText.Store(name, text)
}

func GetText(name string, config *config.ClientConfig) (Text, error) {
	if v, ok := registeredText.Load(name); ok {
		text, err := v.(TextConstructor)(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create text generator: %w", err)
		}
		return text, nil
	}
	return nil, fmt.Errorf("text generator not found: %s", name)
}

type Embeddings interface {
	Embed(ctx context.Context, prompt string) ([]float64, error)
	Size() (uint64, error)
}

var registeredEmbeddings = sync.Map{}

type EmbeddingsConstructor func(*config.ClientConfig) (Embeddings, error)

func RegisterEmbeddings(name string, embeddings EmbeddingsConstructor) {
	registeredEmbeddings.Store(name, embeddings)
}

func GetEmbeddings(name string, config *config.ClientConfig) (Embeddings, error) {
	if v, ok := registeredEmbeddings.Load(name); ok {
		embeddings, err := v.(EmbeddingsConstructor)(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create embeddings client: %w", err)
		}
		return embeddings, nil
	}
	return nil, fmt.Errorf("embeddings client not found: %s", name)
}
