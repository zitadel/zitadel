package eventstore

import "time"

//ReadModel is the minimum representation of a View model.
// It implements a basic reducer
// it might be saved in a database or in memory
type ReadModel struct {
	AggregateID       string        `json:"-"`
	ProcessedSequence uint64        `json:"-"`
	CreationDate      time.Time     `json:"-"`
	ChangeDate        time.Time     `json:"-"`
	Events            []EventReader `json:"-"`
	ResourceOwner     string        `json:"-"`
}

//AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *ReadModel) AppendEvents(events ...EventReader) *ReadModel {
	rm.Events = append(rm.Events, events...)
	return rm
}

//Reduce is the basic implementaion of reducer
// If this function is extended the extending function should be the last step
func (rm *ReadModel) Reduce() error {
	if len(rm.Events) == 0 {
		return nil
	}

	if rm.AggregateID == "" {
		rm.AggregateID = rm.Events[0].AggregateID()
	}
	if rm.ResourceOwner == "" {
		rm.ResourceOwner = rm.Events[0].ResourceOwner()
	}

	if rm.CreationDate.IsZero() {
		rm.CreationDate = rm.Events[0].CreationDate()
	}
	rm.ChangeDate = rm.Events[len(rm.Events)-1].CreationDate()
	rm.ProcessedSequence = rm.Events[len(rm.Events)-1].Sequence()
	// all events processed and not needed anymore
	rm.Events = nil
	rm.Events = []EventReader{}
	return nil
}
