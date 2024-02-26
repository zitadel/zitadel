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
	return eventstore.MergeFilters(
		rm.filter,
		rm.State.Filter,
		rm.PrimaryDomain.Filter,
	)
}

func (rm *Org) filter() []*eventstore.Filter {
	return []*eventstore.Filter{
		// we don't need the filters of the projections as we filter all events of the read model
		eventstore.NewFilter(
			eventstore.AppendAggregateFilter(
				org.AggregateType,
				eventstore.WithAggregateID(rm.ID),
			),
		),
	}
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
