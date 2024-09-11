package projection

import (
	"github.com/zitadel/logging"

	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type InstanceDomain struct {
	Projection

	Domain
}

func NewInstanceDomainFromEvent(event *v2_es.StorageEvent) *InstanceDomain {
	domain := &InstanceDomain{Domain: Domain{id: event.Aggregate.ID}}
	err := domain.reduceAdded(event)
	logging.OnError(err).Error("could not reduce added event")
	return domain
}

func (d *InstanceDomain) reduceAdded(event *v2_es.StorageEvent) error {
	if !d.shouldReduce(event, "") {
		return nil
	}

	e, err := instance.DomainAddedEventFromStorage(event)
	if err != nil {
		return err
	}
	d.reduceAddedPayload(e.Payload)
	d.Projection.Set(event)
	return nil
}

func (d *InstanceDomain) reduceVerified(event *v2_es.StorageEvent) error {
	e, err := instance.DomainVerifiedEventFromStorage(event)
	if err != nil {
		return err
	}
	if !d.shouldReduce(event, e.Payload.Name) {
		return nil
	}

	d.reduceVerifiedPayload(e.Payload)
	d.Projection.Set(event)
	return nil
}

func (d *InstanceDomain) reducePrimarySet(event *v2_es.StorageEvent) error {
	e, err := instance.DomainPrimarySetEventFromStorage(event)
	if err != nil {
		return err
	}
	d.reducePrimarySetPayload(e.Payload)
	d.Projection.Set(event)
	return nil
}
