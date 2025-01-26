package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
)

type Jobs struct {
	pool  *ants.Pool
	errCh chan error
}

func New(ctx context.Context, concurrentWorkerSize int) (*Jobs, error) {
	pool, err := ants.NewPool(concurrentWorkerSize, ants.WithExpiryDuration(3*time.Minute), ants.WithNonblocking(false), ants.WithPanicHandler(func(i interface{}) {
		log.Error().Any("panic", i).Msg("panic in job")
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to create ants pool: %w", err)
	}

	errCh := make(chan error, 1024)

	context.AfterFunc(ctx, func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Any("panic", r).Msg("panic in job closing")
			}
		}()
		pool.Release()
		close(errCh)
	})

	return &Jobs{
		pool:  pool,
		errCh: errCh,
	}, nil
}

func (j *Jobs) Submit(job func() error) <-chan error {
	if err := j.pool.Submit(func() {
		j.errCh <- job()
	}); err != nil {
		j.errCh <- err
	}

	return j.errCh
}

func (j *Jobs) MarkEnd() {
	close(j.errCh)
}

func (j *Jobs) Close() error {
	joinedErr := error(nil)
	for err := range j.errCh {
		if err != nil {
			if joinedErr == nil {
				joinedErr = err
			} else {
				joinedErr = fmt.Errorf("%w; %w", joinedErr, err)
			}
		}
	}

	return joinedErr
}

func (j *Jobs) Err() <-chan error {
	return j.errCh
}
