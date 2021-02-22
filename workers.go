package workers

import (
	"errors"

	"github.com/edge/atomicstore"
)

var (
	ErrDuplicateKey         = errors.New("A worker with this key already exists")
	ErrBroadcastInvalidType = errors.New("Failed to broadcast to Invalid worker")
)

type Workers struct {
	*atomicstore.Store // stores Worker
}

func (w *Workers) Insert(wkr Worker) (Worker, error) {
	// If the worker already exists, return an error and the existing worker.
	if worker, loaded := w.Upsert(wkr.GetID(), wrk); loaded {
		return worker.(*Worker), ErrDuplicateKey
	}

	worker.OnClose(func() {
		w.Delete(wrk.GetID())
	})
	return wrk, nil
}

func (w *Workers) Broadcast(job interface{}) (err error) {
	w.Range(func(k, v interface{}) bool {
		worker, ok := v.(*Worker)
		if !ok {
			err = ErrBroadcastInvalidType
			return true
		}
		worker.JobChan() <- job
		return true
	})
	return
}
