package readmodel

import (
	"context"
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

func (rm *Org) Filter(ctx context.Context) *eventstore.Filter {
	return eventstore.MergeFilters(
		rm.State.Filter(ctx),
		rm.PrimaryDomain.Filter(ctx),
		eventstore.NewFilter(
			ctx,
			eventstore.FilterEventQuery(
				eventstore.FilterAggregateTypes(org.AggregateType),
				eventstore.FilterAggregateIDs(rm.ID),
				eventstore.FilterEventTypes(), // filter for all event types
			),
		),
	)
}

func (rm *Org) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		switch event.Type() {
		case org.Added.Type():
			added := new(org.AddedEvent)
			if err := event.Unmarshal(added); err != nil {
				return err
			}
			rm.Name = added.Name
			rm.Owner = event.Aggregate().Owner
			rm.CreationDate = event.CreatedAt()
		case org.Changed.Type():
			changed := new(org.ChangedEvent)
			if err := event.Unmarshal(changed); err != nil {
				return err
			}
			if changed.Name != nil {
				rm.Name = *changed.Name
			}
		}
		rm.Sequence = event.Sequence()
		rm.ChangeDate = event.CreatedAt()
	}
	if err := rm.State.Reduce(events...); err != nil {
		return err
	}
	return rm.PrimaryDomain.Reduce(events...)
}
