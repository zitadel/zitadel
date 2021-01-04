package eventstore

type aggregater interface {
	//ID returns the aggreagte id
	ID() string
	//KeyType returns the aggregate type
	Type() AggregateType
	//Events returns the events which will be pushed
	Events() []EventPusher
	//ResourceOwner returns the organisation id which manages this aggregate
	// resource owner is only on the inital push needed
	// afterwards the resource owner of the previous event is taken
	ResourceOwner() string
	//Version represents the semantic version of the aggregate
	Version() Version
	//PreviouseSequence should return the sequence of the latest event of this aggregate
	// stored in the eventstore
	// it's set to the first event of this push transaction,
	// later events consume the sequence of the previously pushed event of the aggregate
	PreviousSequence() uint64
}

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

func AggregateFromWriteModel(
	wm *WriteModel,
	typ AggregateType,
	version Version,
) *Aggregate {
	return &Aggregate{
		id:               wm.AggregateID,
		typ:              typ,
		resourceOwner:    wm.ResourceOwner,
		version:          version,
		previousSequence: wm.ProcessedSequence,
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

//KeyType implements aggregater
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
