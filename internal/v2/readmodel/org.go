package readmodel

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/v2/projection"
)

type Org struct {
	Name   string              `json:"name,omitempty"`
	Domain string              `json:"domain,omitempty"`
	State  projection.OrgState `json:"-"`
}

func (rm *Org) Filter() *eventstore.SearchQueryBuilder {
	// TODO: merge search queries from projections
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)
}

func (rm *Org) Reduce(events ...eventstore.Event) error {
	for _, event := range events {
		switch event.Type() {
		case eventstore.EventType(org.Added.Type()):
			added := new(org.AddedEvent)
			event.Unmarshal(added)
			_ = added.Name
		case eventstore.EventType(org.Changed.Type()):
			if err := event.Unmarshal(rm); err != nil {
				return err
			}
		}
	}
	return rm.State.Reduce(events...)
}
