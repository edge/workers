package workers

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrJobHandlerExists = errors.New("This workers job handler already exists")
	ErrInvalidContext   = errors.New("Cannot start worker with invalid context")
)

// DefaultWorker uses the Worker interface.
type DefaultWorker struct {
	sync.RWMutex
	id           string
	jobChan      chan interface{}
	closeHandler func()
	jobHandler   func(interface{})
	metadata     interface{}
	ctx          context.Context
	cancel       context.CancelFunc
}

func (d *DefaultWorker) handleJobs() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case job := <-d.jobChan:
			if d.jobHandler != nil {
				d.jobHandler(job)
			}
		}
	}
}

// GetID returns the ID of the worker.
func (d *DefaultWorker) GetID() string {
	return d.id
}

// Context returns the workers context.
func (d *DefaultWorker) Context() context.Context {
	return d.ctx
}

// Start starts the worker with a context.
func (d *DefaultWorker) Start(ctx context.Context) error {
	// Invalid context
	if ctx == nil || ctx.Err() != nil {
		return ErrInvalidContext
	}

	// Set context and job channel.
	d.ctx, d.cancel = context.WithCancel(ctx)
	d.jobChan = make(chan interface{})

	go d.handleJobs()
	return nil
}

// Stop stops the worker and closes its job channels.
func (d *DefaultWorker) Stop() {
	if d.ctx != nil && d.ctx.Err() == nil {
		d.cancel()
	}
	if d.closeHandler != nil {
		d.closeHandler()
	}
}

// AddJob adds a job to the workers queue.
func (d *DefaultWorker) AddJob(j interface{}) error {
	if d.ctx.Err() != nil {
		return d.ctx.Err()
	}
	d.jobChan <- j
	return nil
}

// ScheduleJob schedules a job for a later time.
func (d *DefaultWorker) ScheduleJob(j interface{}, t time.Time) {
	go func() {
		select {
		case <-d.ctx.Done():
			return
		case <-GetScheduledTicker(t):
			d.AddJob(j)
		}
	}()
}

// JobHandler adds a callback handler for new jobs.
func (d *DefaultWorker) JobHandler(jh func(interface{})) error {
	if d.jobHandler != nil {
		return ErrJobHandlerExists
	}
	d.jobHandler = jh
	return nil
}

// CloseHandler sets a function to call when the worker closes.
func (d *DefaultWorker) CloseHandler(ch func()) {
	d.closeHandler = ch
}

// SetMetadata sets arbitrary metadata values
func (d *DefaultWorker) SetMetadata(md interface{}) {
	d.Lock()
	defer d.Unlock()
	d.metadata = md
}

// GetMetadata gets metadata
func (d *DefaultWorker) GetMetadata() interface{} {
	d.Lock()
	defer d.Unlock()
	return d.metadata
}

// NewDefaultWorker creates an instance of the default worker.
func NewDefaultWorker(id string) Worker {
	return &DefaultWorker{id: id}
}
