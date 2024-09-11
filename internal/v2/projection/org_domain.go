package projection

import (
	"github.com/zitadel/logging"

	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgDomain struct {
	Projection

	Domain
}

func NewOrgDomainFromEvent(event *v2_es.StorageEvent) *OrgDomain {
	domain := &OrgDomain{Domain: Domain{id: event.Aggregate.ID}}
	err := domain.reduceAdded(event)
	logging.OnError(err).Error("could not reduce added event")
	return domain
}

func (d *OrgDomain) reduceAdded(event *v2_es.StorageEvent) error {
	if !d.shouldReduce(event, "") {
		return nil
	}

	e, err := org.DomainAddedEventFromStorage(event)
	if err != nil {
		return err
	}
	d.reduceAddedPayload(e.Payload)
	d.Projection.Set(event)
	return nil
}

func (d *OrgDomain) reduceVerified(event *v2_es.StorageEvent) error {
	e, err := org.DomainVerifiedEventFromStorage(event)
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

func (d *OrgDomain) reducePrimarySet(event *v2_es.StorageEvent) error {
	e, err := org.DomainPrimarySetEventFromStorage(event)
	if err != nil {
		return err
	}
	d.reducePrimarySetPayload(e.Payload)
	d.Projection.Set(event)
	return nil
}
