package projection

import (
	"github.com/zitadel/logging"

	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type Instance struct {
	ObjectMetadata

	AuthZInstance

	Name string
}

func NewInstanceFromEvent(event *v2_es.StorageEvent) *Instance {
	instance := &Instance{
		AuthZInstance: *NewAuthZInstanceFromEvent(event),
	}
	err := instance.reduceAdded(event)
	logging.OnError(err).Error("could not reduce added event")

	return instance
}

func (i *Instance) Reducers() Reducers {
	if i.ObjectMetadata.Reducers != nil {
		return i.ObjectMetadata.Reducers
	}
	i.ObjectMetadata.Reducers = Reducers{
		instance.AggregateType: {
			instance.AddedType:              i.reduceAdded,
			instance.ChangedType:            i.reduceChanged,
			instance.DefaultOrgSetType:      i.reduce,
			instance.ProjectSetType:         i.reduce,
			instance.ConsoleSetType:         i.reduce,
			instance.DefaultLanguageSetType: i.reduce,
			instance.RemovedType:            i.reduce,

			instance.DomainAddedType:      i.reduce,
			instance.DomainVerifiedType:   i.reduce,
			instance.DomainPrimarySetType: i.reduce,
			instance.DomainRemovedType:    i.reduce,
		},
	}

	return i.ObjectMetadata.Reducers
}

func (i *Instance) reduce(event *v2_es.StorageEvent) error {
	if !i.ObjectMetadata.ShouldReduce(event) {
		return nil
	}
	return i.ObjectMetadata.Reduce(event, i.AuthZInstance.Reducers()[event.Aggregate.Type][event.Type])
}

func (i *Instance) reduceAdded(event *v2_es.StorageEvent) error {
	if !i.ObjectMetadata.ShouldReduce(event) {
		return nil
	}

	e, err := instance.AddedEventFromStorage(event)
	if err != nil {
		return err
	}

	err = i.ObjectMetadata.Reduce(event, i.AuthZInstance.reduceAdded)
	if err != nil {
		return err
	}
	i.Name = e.Payload.Name

	return nil
}

func (i *Instance) reduceChanged(event *v2_es.StorageEvent) error {
	if !i.ObjectMetadata.ShouldReduce(event) {
		return nil
	}

	e, err := instance.ChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	i.Name = e.Payload.Name
	i.ObjectMetadata.Set(event)
	return nil
}
