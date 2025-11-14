package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceTrustedDomainsWriteModel struct {
	eventstore.WriteModel

	Domains []string
}

func NewInstanceTrustedDomainsWriteModel(ctx context.Context) *InstanceTrustedDomainsWriteModel {
	instanceID := authz.GetInstance(ctx).InstanceID()
	return &InstanceTrustedDomainsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
	}
}

func (wm *InstanceTrustedDomainsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.TrustedDomainAddedEvent:
			wm.Domains = append(wm.Domains, e.Domain)
		case *instance.TrustedDomainRemovedEvent:
			wm.Domains = slices.DeleteFunc(wm.Domains, func(domain string) bool {
				return domain == e.Domain
			})
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceTrustedDomainsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.TrustedDomainAddedEventType,
			instance.TrustedDomainRemovedEventType,
		).
		Builder()
}
