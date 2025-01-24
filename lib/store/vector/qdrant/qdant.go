package qdrant

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/qdrant/go-client/qdrant"

	"github.com/snowmerak/indexer/pkg/box"
)

type Config struct {
	host       string
	port       int
	name       string
	apiKey     string
	useTLS     bool
	volatility bool
}

func NewConfig(host string, port int, name string) *Config {
	return &Config{
		host: host,
		port: port,
		name: name,
	}
}

func (c *Config) WithAPIKey(apiKey string) *Config {
	c.apiKey = apiKey
	return c
}

func (c *Config) WithTLS() *Config {
	c.useTLS = true
	return c
}

func (c *Config) WithVolatility() *Config {
	c.volatility = true
	return c
}

type Vector struct {
	client *qdrant.Client
	config *Config
	id     atomic.Int64
}

func New(ctx context.Context, cfg *Config) (*Vector, error) {
	qc := &qdrant.Config{
		Host: cfg.host,
		Port: cfg.port,
	}

	if cfg.apiKey != "" {
		qc.APIKey = cfg.apiKey
	}

	if cfg.useTLS {
		qc.UseTLS = true
	}

	client, err := qdrant.NewClient(qc)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		client.Close()
	})

	if err := client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: cfg.name,
	}); err != nil {
		if cfg.volatility {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	if cfg.volatility {
		context.AfterFunc(ctx, func() {
			count := 0
		loop:
			for {
				if err := client.DeleteCollection(context.Background(), cfg.name); err != nil {
					count++
					if count > 3 {
						break loop
					}
					continue loop
				}

				break loop
			}
		})
	}

	return &Vector{
		client: client,
		config: cfg,
	}, nil
}

func (v *Vector) NextId() int {
	return int(v.id.Add(1))
}

func (v *Vector) Store(ctx context.Context, id int, vector []float64) error {
	nv := make([]float32, len(vector))
	for i, v := range vector {
		nv[i] = float32(v)
	}

	if _, err := v.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: v.config.name,
		Wait:           box.Wrap(true),
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(uint64(id)),
				Vectors: qdrant.NewVectors(nv...),
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}

func (v *Vector) Get(ctx context.Context, id int) ([]float64, error) {
	pt, err := v.client.Get(ctx, &qdrant.GetPoints{
		CollectionName: v.config.name,
		Ids:            []*qdrant.PointId{qdrant.NewIDNum(uint64(id))},
		WithVectors:    qdrant.NewWithVectors(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}

	if len(pt) == 0 {
		return nil, fmt.Errorf("point not found")
	}

	fpt := pt[0]
	vector := fpt.GetVectors().GetVector().Data
	nv := make([]float64, len(vector))

	for i, v := range vector {
		nv[i] = float64(v)
	}

	return nv, nil
}

func (v *Vector) Search(ctx context.Context, vector []float64, limit int) ([]int, error) {
	nv := make([]float32, len(vector))
	for i, v := range vector {
		nv[i] = float32(v)
	}

	resp, err := v.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: v.config.name,
		Query:          qdrant.NewQuery(nv...),
		Limit:          box.Wrap(uint64(limit)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	ids := make([]int, len(resp))
	for i, pt := range resp {
		ids[i] = int(pt.GetId().GetNum())
	}

	return ids, nil
}

func (v *Vector) Delete(ctx context.Context, id int) error {
	if _, err := v.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: v.config.name,
		Points:         qdrant.NewPointsSelector(qdrant.NewIDNum(uint64(id))),
	}); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

func (v *Vector) Drop(ctx context.Context) error {
	if err := v.client.DeleteCollection(ctx, v.config.name); err != nil {
		return fmt.Errorf("failed to drop: %w", err)
	}

	return nil
}
