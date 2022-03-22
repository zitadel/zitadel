package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type InstanceDomainWriteModel struct {
	eventstore.WriteModel

	Domain    string
	Generated bool
	State     domain.InstanceDomainState
}

func NewInstanceDomainWriteModel(instanceDomain string) *InstanceDomainWriteModel {
	return &InstanceDomainWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
		Domain: instanceDomain,
	}
}

func (wm *InstanceDomainWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.DomainAddedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *iam.DomainRemovedEvent:
			if e.Domain != wm.Domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceDomainWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.DomainAddedEvent:
			wm.Domain = e.Domain
			wm.Generated = e.Generated
			wm.State = domain.InstanceDomainStateActive
		case *iam.DomainRemovedEvent:
			wm.State = domain.InstanceDomainStateRemoved
		}
	}
	return nil
}

func (wm *InstanceDomainWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(domain.IAMID).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(domain.IAMID).
		EventTypes(
			iam.InstanceDomainAddedEventType,
			iam.InstanceDomainRemovedEventType).
		Builder()
}
