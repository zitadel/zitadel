package eventstore

//ReadModel is the minimum representation of a View model.
// it might be saved in a database or in memory
type ReadModel struct {
	ProcessedSequence uint64
	ID                string
	Events            []Event
}

//Append adds all the events to the aggregate.
// The function doesn't compute the new state of the read model
func (a *ReadModel) Append(events ...Event) {
	a.Events = append(a.Events, events...)
}
