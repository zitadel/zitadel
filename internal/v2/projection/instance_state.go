package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type InstanceState struct {
	projection

	id string

	instance.State
}

func NewInstanceStateProjection(id string) *InstanceState {
	return &InstanceState{
		id: id,
	}
}

func (p *InstanceState) Reduce(events ...*eventstore.StorageEvent) error {
	for _, event := range events {
		if !p.shouldReduce(event) {
			continue
		}

		switch event.Type {
		case instance.AddedType:
			p.State = instance.ActiveState
		case instance.RemovedType:
			p.State = instance.RemovedState
		default:
			continue
		}
		p.position = event.Position
	}
	return nil
}
