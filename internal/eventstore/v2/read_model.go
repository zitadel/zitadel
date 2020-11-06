package eventstore

import "time"

func NewReadModel( /*id string*/ ) *ReadModel {
	return &ReadModel{
		// ID:     id,
		Events: []EventReader{},
	}
}

//ReadModel is the minimum representation of a View model.
// it might be saved in a database or in memory
type ReadModel struct {
	ProcessedSequence uint64 `json:"-"`
	// ID                string        `json:"-"`
	CreationDate time.Time     `json:"-"`
	ChangeDate   time.Time     `json:"-"`
	Events       []EventReader `json:"-"`
}

//AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *ReadModel) AppendEvents(events ...EventReader) *ReadModel {
	rm.Events = append(rm.Events, events...)
	return rm
}

//Reduce must be the last step in the reduce function of the extension
func (rm *ReadModel) Reduce() error {
	if len(rm.Events) == 0 {
		return nil
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
