package eventstore

func NewReadModel() *ReadModel {
	return &ReadModel{
		Events: []EventReader{},
	}
}

//WriteModel is the minimum representation of a command side view model.
// It implements a basic reducer
// it's purpose is to reduce events to create new ones
type WriteModel struct {
	AggregateID       string        `json:"-"`
	ProcessedSequence uint64        `json:"-"`
	Events            []EventReader `json:"-"`
}

//AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *WriteModel) AppendEvents(events ...EventReader) *WriteModel {
	rm.Events = append(rm.Events, events...)
	return rm
}

//Reduce is the basic implementaion of reducer
// If this function is extended the extending function should be the last step
func (rm *WriteModel) Reduce() error {
	if len(rm.Events) == 0 {
		return nil
	}

	if rm.AggregateID == "" {
		rm.AggregateID = rm.Events[0].AggregateID()
	}
	if rm.ResourceOwner == "" {
		rm.ResourceOwner = rm.Events[0].ResourceOwner()
	}

	rm.ProcessedSequence = rm.Events[len(rm.Events)-1].Sequence()

	// all events processed and not needed anymore
	rm.Events = nil
	rm.Events = []EventReader{}
	return nil
}
