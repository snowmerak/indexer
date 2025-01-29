package vector

import (
	"context"
	"fmt"
	"sync"

	"github.com/snowmerak/indexer/pkg/config"
)

const (
	PayloadId        = "id"
	PayloadVector    = "vector"
	PayloadRelatedId = "related_id"
)

type Payload struct {
	Id        int
	Vector    []float64
	RelatedId int
}

type Vector interface {
	Create(ctx context.Context, vectorSize uint64) error
	Store(ctx context.Context, id int, payload *Payload) error
	Get(ctx context.Context, id int) (*Payload, error)
	Search(ctx context.Context, vector []float64, limit int) ([]*Payload, error)
	Delete(ctx context.Context, id int) error
	Drop(ctx context.Context) error
}

var registeredVector = sync.Map{}

type VectorConstructor func(*config.ClientConfig) (Vector, error)

func RegisterVector(name string, vector VectorConstructor) {
	registeredVector.Store(name, vector)
}

func GetVector(name string, config *config.ClientConfig) (Vector, error) {
	if v, ok := registeredVector.Load(name); ok {
		vector, err := v.(VectorConstructor)(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create vector client: %w", err)
		}
		return vector, nil
	}
	return nil, fmt.Errorf("vector client not found: %s", name)
}
