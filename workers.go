package workers

import (
	"context"
	"errors"
	"fmt"

	"github.com/edge/atomicstore"
)

var (
	ErrBroadcastInvalidType = errors.New("Failed to broadcast to Invalid worker")
)

type Pool struct {
	*atomicstore.Store // stores Worker
}

func (p *Pool) Add(ctx context.Context, wrk Worker) (Worker, error) {
	// If the worker already exists, return an error and the existing worker.
	if worker, loaded := p.InsertUnique(wrk.GetID(), wrk); loaded {
		return worker.(Worker), fmt.Errorf(`A worker with key %s already exists`, wrk.GetID())
	}

	err := wrk.Start(ctx)

	return wrk, err
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
		Store: atomicstore.New(false),
	}
}
