package eventstore

import (
	"time"
)

// WriteModel is the minimum representation of a command side write model.
// It implements a basic reducer
// it's purpose is to reduce events to create new ones
type WriteModel struct {
	AggregateID       string    `json:"-"`
	ProcessedSequence uint64    `json:"-"`
	Events            []Event   `json:"-"`
	ResourceOwner     string    `json:"-"`
	InstanceID        string    `json:"-"`
	ChangeDate        time.Time `json:"-"`
}

// AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *WriteModel) AppendEvents(events ...Event) {
	rm.Events = append(rm.Events, events...)
}

// Reduce is the basic implementation of reducer
// If this function is extended the extending function should be the last step
func (wm *WriteModel) Reduce() error {
	if len(wm.Events) == 0 {
		return nil
	}

	latestEvent := wm.Events[len(wm.Events)-1]
	if wm.AggregateID == "" {
		wm.AggregateID = latestEvent.Aggregate().ID
	}

	if wm.ResourceOwner == "" {
		wm.ResourceOwner = latestEvent.Aggregate().ResourceOwner
	}

	if wm.InstanceID == "" {
		wm.InstanceID = latestEvent.Aggregate().InstanceID
	}

	wm.ProcessedSequence = latestEvent.Sequence()
	wm.ChangeDate = latestEvent.CreatedAt()

	// all events processed and not needed anymore
	wm.Events = nil
	wm.Events = []Event{}

	return nil
}
