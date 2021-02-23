package workers

import (
	"context"
	"time"
)

// Worker is a contributing entry to the worker Pool.
type Worker interface {
	// Start the worker with the given context
	GetID() string
	// Context returns the workers context
	Context() context.Context
	// Start the worker with the given context
	Start(context.Context) error
	// Stop the worker
	Stop()
	// AddJob to the workers queue
	AddJob(interface{}) error
	// ScheduleJob
	ScheduleJob(interface{}, time.Time)
	// JobHandler
	JobHandler(func(interface{})) error
	// CloseHandler
	CloseHandler(func())
	// SetMetadata sets arbitrary metadata values
	SetMetadata(interface{})
	// GetMetadata gets metadata
	GetMetadata() interface{}
}
