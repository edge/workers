package workers

import (
	"context"
	"errors"

	"github.com/edge/atomicstore"
)

var (
	ErrDuplicateKey         = errors.New("A worker with this key already exists")
	ErrBroadcastInvalidType = errors.New("Failed to broadcast to Invalid worker")
)

type Pool struct {
	*atomicstore.Store // stores Worker
}

func (p *Pool) Add(ctx context.Context, wrk Worker) (Worker, error) {
	// If the worker already exists, return an error and the existing worker.
	if worker, loaded := p.Upsert(wrk.GetID(), wrk); loaded {
		return worker.(Worker), ErrDuplicateKey
	}

	wrk.Start(ctx)

	return wrk, nil
}

// Remove stops the worker and removes it from the pool.
func (p *Pool) Remove(key string) {
	if w, exists := p.Get(key); exists {
		worker := w.(Worker)
		worker.Stop()
		p.Delete(key)
	}
}

func (p *Pool) Broadcast(job interface{}) (err error) {
	p.Range(func(k, v interface{}) bool {
		worker, ok := v.(Worker)
		if !ok {
			err = ErrBroadcastInvalidType
			return true
		}
		worker.AddJob(job)
		return true
	})
	return
}

// New returns a new instance of workers.
func New() *Pool {
	return &Pool{
		Store: atomicstore.New(true),
	}
}
