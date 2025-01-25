package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
)

type Jobs struct {
	pool *ants.Pool
}

func New(ctx context.Context, size int) (*Jobs, error) {
	pool, err := ants.NewPool(size, ants.WithExpiryDuration(3*time.Minute), ants.WithNonblocking(false), ants.WithPanicHandler(func(i interface{}) {
		log.Error().Any("panic", i).Msg("panic in job")
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to create ants pool: %w", err)
	}

	context.AfterFunc(ctx, func() {
		pool.Release()
	})

	return &Jobs{
		pool: pool,
	}, nil
}

func (j *Jobs) Submit(job func() error) <-chan error {
	ch := make(chan error, 1)

	if err := j.pool.Submit(func() {
		ch <- job()
	}); err != nil {
		close(ch)
	}

	return ch
}
