package projection

import (
	"github.com/zitadel/logging"

	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type Instance struct {
	objectMetadata

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

func (i *Instance) Reducers() map[string]map[string]v2_es.ReduceEvent {
	if i.objectMetadata.reducers != nil {
		return i.objectMetadata.reducers
	}
	i.objectMetadata.reducers = map[string]map[string]v2_es.ReduceEvent{
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

	return i.objectMetadata.reducers
}

func (i *Instance) reduce(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}
	return i.objectMetadata.reduce(event, i.AuthZInstance.Reducers()[event.Aggregate.Type][event.Type])
}

func (i *Instance) reduceAdded(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.AddedEventFromStorage(event)
	if err != nil {
		return err
	}

	err = i.objectMetadata.reduce(event, i.AuthZInstance.reduceAdded)
	if err != nil {
		return err
	}
	i.Name = e.Payload.Name

	return nil
}

func (i *Instance) reduceChanged(event *v2_es.StorageEvent) error {
	if !i.ShouldReduce(event) {
		return nil
	}

	e, err := instance.ChangedEventFromStorage(event)
	if err != nil {
		return err
	}
	i.Name = e.Payload.Name
	i.objectMetadata.set(event)
	return nil
}
