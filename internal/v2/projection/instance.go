package projection

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type Instance struct {
	AuthZInstance

	Name string

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint32
}

// InterestedIn implements model.
func (i *Instance) InterestedIn() map[eventstore.AggregateType][]eventstore.EventType {
	return map[eventstore.AggregateType][]eventstore.EventType{
		eventstore.AggregateType(instance.AggregateType): {
			eventstore.EventType(instance.AddedType),
			eventstore.EventType(instance.ChangedType),
			eventstore.EventType(instance.DefaultOrgSetType),
			eventstore.EventType(instance.ProjectSetType),
			eventstore.EventType(instance.ConsoleSetType),
			eventstore.EventType(instance.DefaultLanguageSetType),
			eventstore.EventType(instance.RemovedType),
		},
	}
}

func (i *Instance) Reduce(events ...*v2_es.StorageEvent) (err error) {
	for _, event := range events {
		if err = i.AuthZInstance.Reduce(event); err != nil {
			return err
		}

		switch event.Type {
		case instance.AddedType:
			err = i.reduceAdded(event)
		case instance.ChangedType:
			err = i.reduceChanged(event)
		}
		if err != nil {
			return err
		}

		if err = i.State.Reduce(event); err != nil {
			return err
		}

		i.ChangeDate = event.CreatedAt
		i.Sequence = event.Sequence
	}
	return nil
}

func (i *Instance) reduceAdded(event *v2_es.StorageEvent) error {
	e, err := instance.AddedEventFromStorage(event)
	if err != nil {
		return err
	}
	i.CreationDate = e.CreatedAt
	i.Name = e.Payload.Name
	return nil
}

func (i *Instance) reduceChanged(event *v2_es.StorageEvent) error {
	e, err := instance.ChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	i.Name = e.Payload.Name
	return nil
}
