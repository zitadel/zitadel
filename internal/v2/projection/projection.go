package projection

import "github.com/zitadel/zitadel/internal/v2/eventstore"

type projection struct {
	instance string
	position eventstore.GlobalPosition
}

func (p *projection) reduce(event *eventstore.Event[eventstore.StoragePayload]) {
	if p.instance == "" {
		p.instance = event.Aggregate.Instance
	}
	p.position = event.Position
}

func (p projection) shouldReduce(event *eventstore.Event[eventstore.StoragePayload]) bool {
	shouldReduce := p.instance == "" || p.instance == event.Aggregate.Instance
	if p.position.Position == event.Position.Position {
		shouldReduce = shouldReduce && p.position.InPositionOrder < event.Position.InPositionOrder
	} else {
		shouldReduce = shouldReduce && p.position.Position < event.Position.Position
	}

	return shouldReduce
}
