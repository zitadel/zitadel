package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceAllowedDomainsWriteModel struct {
	eventstore.WriteModel

	Domains []string
}

func NewInstanceAllowedDomainsWriteModel(ctx context.Context) *InstanceAllowedDomainsWriteModel {
	return &InstanceAllowedDomainsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		},
	}
}

func (wm *InstanceAllowedDomainsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.AllowedDomainAddedEvent:
			wm.Domains = append(wm.Domains, e.Domain)
		case *instance.AllowedDomainRemovedEvent:
			wm.Domains = slices.DeleteFunc(wm.Domains, func(domain string) bool {
				return domain == e.Domain
			})
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceAllowedDomainsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.AllowedDomainAddedEventType,
			instance.AllowedDomainRemovedEventType,
		).
		Builder()
}
