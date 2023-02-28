package eventstore

import "time"

// ReadModel is the minimum representation of a read model.
// It implements a basic reducer
// it might be saved in a database or in memory
type ReadModel struct {
	AggregateID       string    `json:"-"`
	ProcessedSequence uint64    `json:"-"`
	CreationDate      time.Time `json:"-"`
	ChangeDate        time.Time `json:"-"`
	Events            []Event   `json:"-"`
	ResourceOwner     string    `json:"-"`
	InstanceID        string    `json:"-"`
}

// AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *ReadModel) AppendEvents(events ...Event) {
	rm.Events = append(rm.Events, events...)
}

// Reduce is the basic implementation of reducer
// If this function is extended the extending function should be the last step
func (rm *ReadModel) Reduce() error {
	if len(rm.Events) == 0 {
		return nil
	}

	if rm.AggregateID == "" {
		rm.AggregateID = rm.Events[0].Aggregate().ID
	}
	if rm.ResourceOwner == "" {
		rm.ResourceOwner = rm.Events[0].Aggregate().ResourceOwner
	}
	if rm.InstanceID == "" {
		rm.InstanceID = rm.Events[0].Aggregate().InstanceID
	}

	if rm.CreationDate.IsZero() {
		rm.CreationDate = rm.Events[0].CreationDate()
	}
	rm.ChangeDate = rm.Events[len(rm.Events)-1].CreationDate()
	rm.ProcessedSequence = rm.Events[len(rm.Events)-1].Sequence()
	// all events processed and not needed anymore
	rm.Events = nil
	rm.Events = []Event{}
	return nil
}
