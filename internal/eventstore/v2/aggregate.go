package eventstore

func NewAggregate(
	id string,
	typ AggregateType,
	resourceOwner string,
	version Version,
	previousSequence uint64,
) *Aggregate {
	return &Aggregate{
		id:               id,
		typ:              typ,
		resourceOwner:    resourceOwner,
		version:          version,
		previousSequence: previousSequence,
		events:           []EventPusher{},
	}
}

//Aggregate is the basic implementation of aggregater
type Aggregate struct {
	id               string        `json:"-"`
	typ              AggregateType `json:"-"`
	events           []EventPusher `json:"-"`
	resourceOwner    string        `json:"-"`
	version          Version       `json:"-"`
	previousSequence uint64        `json:"-"`
}

//PushEvents adds all the events to the aggregate.
// The added events will be pushed to eventstore
func (a *Aggregate) PushEvents(events ...EventPusher) *Aggregate {
	a.events = append(a.events, events...)
	return a
}

//ID implements aggregater
func (a *Aggregate) ID() string {
	return a.id
}

//Type implements aggregater
func (a *Aggregate) Type() AggregateType {
	return a.typ
}

//Events implements aggregater
func (a *Aggregate) Events() []EventPusher {
	return a.events
}

//ResourceOwner implements aggregater
func (a *Aggregate) ResourceOwner() string {
	return a.resourceOwner
}

//Version implements aggregater
func (a *Aggregate) Version() Version {
	return a.version
}

//PreviousSequence implements aggregater
func (a *Aggregate) PreviousSequence() uint64 {
	return a.previousSequence
}
