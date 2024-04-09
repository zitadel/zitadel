package readmodel

import (
	"time"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/v2/projection"
)

type Org struct {
	ID            string
	Name          string
	PrimaryDomain *projection.OrgPrimaryDomain
	State         *projection.OrgState

	Sequence     uint32
	CreationDate time.Time
	ChangeDate   time.Time
	Owner        string
}

func NewOrg(id string) *Org {
	return &Org{
		ID:            id,
		State:         projection.NewStateProjection(id),
		PrimaryDomain: projection.NewOrgPrimaryDomain(id),
	}
}

func (rm *Org) Filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		// we don't need the filters of the projections as we filter all events of the read model
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.AggregateID(rm.ID),
			),
		),
	}
}

func (rm *Org) Reduce(events ...*eventstore.Event[eventstore.StoragePayload]) error {
	for _, event := range events {
		switch {
		case org.Added.IsType(event.Type):
			added, err := org.AddedEventFromStorage(event)
			if err != nil {
				return err
			}
			rm.Name = added.Payload.Name
			rm.Owner = event.Aggregate.Owner
			rm.CreationDate = event.CreatedAt
		case org.Changed.IsType(event.Type):
			changed, err := org.ChangedEventFromStorage(event)
			if err != nil {
				return err
			}
			if changed.Payload.Name != nil {
				rm.Name = *changed.Payload.Name
			}
		}
		rm.Sequence = event.Sequence
		rm.ChangeDate = event.CreatedAt
	}
	if err := rm.State.Reduce(events...); err != nil {
		return err
	}
	return rm.PrimaryDomain.Reduce(events...)
}
