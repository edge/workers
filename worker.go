package Workers

import (
	"context"
	"errors"
)

type Worker interface {
	GetID() string
	DoJob(interface{}) error
	StartWithContext(context.Context) (bool, error)
	OnClose(func())
}

var (
	ErrInvalidContext = errors.New("Cannot start worker with invalid context")
	ErrNotRunning     = errors.New("Cannot execute task on stopped worker")
	ErrAlreadyRunning = errors.New("Cannot start a running worker")
)

// DefaultWorker uses the Worker interface.
type DefaultWorker struct {
	id      string
	jobChan chan interface{}
	onClose func()
	running bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// GetID returns the ID of the worker.
func (d *DefaultWorker) GetID() interface{} {
	return d.id
}

// DoJob executes a job if the workers context has no errors.
func (d *DefaultWorker) DoJob(j interface{}) error {
	if !d.running {
		return ErrNotRunning
	}
	if d.ctx.Err() != nil {
		return d.ctx.Err()
	}
	d.jobChan <- j
	return nil
}

// Start sets the worker as started and launches a go routine to stop the worker when its context closes.
func (d *DefaultWorker) Start() (bool, error) {
	return d.StartWithContext(context.Background())
}

// StartWithContext starts the worker with a context.
func (d *DefaultWorker) StartWithContext(ctx context.Context) (bool, error) {
	if d.running {
		return false, ErrAlreadyRunning
	}

	if ctx == nil || ctx.Err() != nil {
		return false, ErrInvalidContext
	}

	go func() {
		<-ctx.Done()
		close(d.jobChan)
		if d.onClose != nil {
			d.onClose()
		}
		d.jobChan = nil
		d.running = false
	}()

	d.running = true
	return true, nil
}

// OnClose sets a function to call when the worker closes.
func (d *DefaultWorker) OnClose(onClose func()) {
	d.onClose = onClose
}
