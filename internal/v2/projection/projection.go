package projection

import "github.com/zitadel/zitadel/internal/v2/eventstore"

type projection struct {
	instance string
	position eventstore.GlobalPosition
}

func (p *projection) reduce(event *eventstore.StorageEvent) {
	if p.instance == "" {
		p.instance = event.Aggregate.Instance
	}
	p.position = event.Position
}

func (p *projection) shouldReduce(event *eventstore.StorageEvent) bool {
	shouldReduce := p.instance == "" || p.instance == event.Aggregate.Instance
	return shouldReduce && p.position.IsLess(event.Position)
}
