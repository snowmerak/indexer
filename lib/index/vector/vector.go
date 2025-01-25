package vector

import "context"

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
	Create(ctx context.Context) error
	Store(ctx context.Context, id int, payload *Payload) error
	Get(ctx context.Context, id int) (*Payload, error)
	Search(ctx context.Context, vector []float64, limit int) ([]*Payload, error)
	Delete(ctx context.Context, id int) error
	Drop(ctx context.Context) error
}
