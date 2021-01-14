package eventstore

type Aggregater interface {
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
}

func NewAggregate(
	id string,
	typ AggregateType,
	resourceOwner string,
	version Version,
) *Aggregate {
	return &Aggregate{
		id:            id,
		typ:           typ,
		resourceOwner: resourceOwner,
		version:       version,
		events:        []EventPusher{},
	}
}

func AggregateFromWriteModel(
	wm *WriteModel,
	typ AggregateType,
	version Version,
) *Aggregate {
	return &Aggregate{
		id:            wm.AggregateID,
		typ:           typ,
		resourceOwner: wm.ResourceOwner,
		version:       version,
		events:        []EventPusher{},
	}
}

//Aggregate is the basic implementation of Aggregater
type Aggregate struct {
	id            string        `json:"-"`
	typ           AggregateType `json:"-"`
	events        []EventPusher `json:"-"`
	resourceOwner string        `json:"-"`
	version       Version       `json:"-"`
}

//PushEvents adds all the events to the aggregate.
// The added events will be pushed to eventstore
func (a *Aggregate) PushEvents(events ...EventPusher) *Aggregate {
	a.events = append(a.events, events...)
	return a
}

//ID implements Aggregater
func (a *Aggregate) ID() string {
	return a.id
}

//KeyType implements Aggregater
func (a *Aggregate) Type() AggregateType {
	return a.typ
}

//Events implements Aggregater
func (a *Aggregate) Events() []EventPusher {
	return a.events
}

//ResourceOwner implements Aggregater
func (a *Aggregate) ResourceOwner() string {
	return a.resourceOwner
}

//Version implements Aggregater
func (a *Aggregate) Version() Version {
	return a.version
}
