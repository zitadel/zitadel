package eventstore

import "time"

func NewReadModel(id string) *ReadModel {
	return &ReadModel{
		ID:     id,
		Events: []Event{},
	}
}

//ReadModel is the minimum representation of a View model.
// it might be saved in a database or in memory
type ReadModel struct {
	ProcessedSequence uint64    `json:"-"`
	ID                string    `json:"-"`
	CreationDate      time.Time `json:"-"`
	ChangeDate        time.Time `json:"-"`
	Events            []Event   `json:"-"`
}

//AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *ReadModel) AppendEvents(events ...Event) *ReadModel {
	rm.Events = append(rm.Events, events...)
	return rm
}

//Reduce must be the last step in the reduce function of the extension
func (rm *ReadModel) Reduce() error {
	if len(rm.Events) == 0 {
		return nil
	}

	if rm.CreationDate.IsZero() {
		rm.CreationDate = rm.Events[0].MetaData().CreationDate
	}
	rm.ChangeDate = rm.Events[len(rm.Events)-1].MetaData().CreationDate
	rm.ProcessedSequence = rm.Events[len(rm.Events)-1].MetaData().Sequence
	// all events processed and not needed anymore
	rm.Events = nil
	rm.Events = []Event{}
	return nil
}

func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		ID:     id,
		Events: []Event{},
	}
}

type Aggregate struct {
	PreviousSequence uint64  `json:"-"`
	ID               string  `json:"-"`
	Events           []Event `json:"-"`
}

//AppendEvents adds all the events to the aggregate.
// The function doesn't compute the new state of the aggregate
func (a *Aggregate) AppendEvents(events ...Event) *Aggregate {
	a.Events = append(a.Events, events...)
	return a
}

//Reduce must be the last step in the reduce function of the extension
func (a *Aggregate) Reduce() error {
	if len(a.Events) == 0 {
		return nil
	}

	a.PreviousSequence = a.Events[len(a.Events)-1].MetaData().Sequence
	// all events processed and not needed anymore
	a.Events = nil
	a.Events = []Event{}
	return nil
}
