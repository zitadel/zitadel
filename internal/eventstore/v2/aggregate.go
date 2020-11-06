package eventstore

func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		ID:     id,
		Events: []EventPusher{},
	}
}

type Aggregate struct {
	PreviousSequence uint64        `json:"-"`
	ID               string        `json:"-"`
	Events           []EventPusher `json:"-"`
}

//AppendEvents adds all the events to the aggregate.
// The function doesn't compute the new state of the aggregate
func (a *Aggregate) AppendEvents(events ...EventPusher) *Aggregate {
	a.Events = append(a.Events, events...)
	return a
}

//Reduce must be the last step in the reduce function of the extension
func (a *Aggregate) Reduce() error {
	if len(a.Events) == 0 {
		return nil
	}

	a.PreviousSequence = a.Events[len(a.Events)-1].Sequence()
	// all events processed and not needed anymore
	a.Events = nil
	a.Events = []Event{}
	return nil
}
