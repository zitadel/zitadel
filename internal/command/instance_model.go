package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"golang.org/x/text/language"
)

type InstanceWriteModel struct {
	eventstore.WriteModel

	Name            string
	State           domain.InstanceState
	GeneratedDomain string

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	GlobalOrgID     string
	ProjectID       string
	DefaultLanguage language.Tag
}

func NewInstanceWriteModel(instanceID string) *InstanceWriteModel {
	return &InstanceWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
	}
}

func (wm *InstanceWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.InstanceAddedEvent:
			wm.Name = e.Name
			wm.State = domain.InstanceStateActive
		case *instance.InstanceChangedEvent:
			wm.Name = e.Name
		case *instance.InstanceRemovedEvent:
			wm.State = domain.InstanceStateRemoved
		case *instance.DomainAddedEvent:
			if !e.Generated {
				continue
			}
			wm.GeneratedDomain = e.Domain
		case *instance.ProjectSetEvent:
			wm.ProjectID = e.ProjectID
		case *instance.GlobalOrgSetEvent:
			wm.GlobalOrgID = e.OrgID
		case *instance.DefaultLanguageSetEvent:
			wm.DefaultLanguage = e.Language
		case *instance.SetupStepEvent:
			if e.Done {
				wm.SetUpDone = e.Step
			} else {
				wm.SetUpStarted = e.Step
			}
		}
	}
	return nil
}

func (wm *InstanceWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.ProjectSetEventType,
			instance.GlobalOrgSetEventType,
			instance.DefaultLanguageSetEventType,
			instance.SetupStartedEventType,
			instance.SetupDoneEventType,
			instance.InstanceAddedEventType,
			instance.InstanceChangedEventType,
			instance.InstanceRemovedEventType,
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
			instance.ProjectSetEventType,
			instance.GlobalOrgSetEventType,
			instance.DefaultLanguageSetEventType,
			instance.SetupStartedEventType,
			instance.SetupDoneEventType).
		Builder()
}

func InstanceAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, instance.AggregateType, instance.AggregateVersion)
}
