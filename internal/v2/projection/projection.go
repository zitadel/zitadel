package projection

import "github.com/zitadel/zitadel/internal/v2/eventstore"

type projection struct {
	instance string
	position eventstore.GlobalPosition
	sequence uint32
}

func (p projection) shouldReduce(event eventstore.Event) bool {
	shouldReduce := p.instance == "" || p.instance == event.Aggregate().Instance
	if p.position.Position == event.Position().Position {
		shouldReduce = shouldReduce && p.position.InPositionOrder < event.Position().InPositionOrder
	} else {
		shouldReduce = shouldReduce && p.position.Position < event.Position().Position
	}

	return shouldReduce
}
