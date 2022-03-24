package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"golang.org/x/text/language"
)

type IAMWriteModel struct {
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

func NewIAMWriteModel(instanceID string) *IAMWriteModel {
	return &IAMWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
	}
}

func (wm *IAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.InstanceAddedEvent:
			wm.Name = e.Name
			wm.State = domain.InstanceStateActive
		case *iam.InstanceChangedEvent:
			wm.Name = e.Name
		case *iam.InstanceRemovedEvent:
			wm.State = domain.InstanceStateRemoved
		case *iam.DomainAddedEvent:
			if !e.Generated {
				continue
			}
			wm.GeneratedDomain = e.Domain
		case *iam.ProjectSetEvent:
			wm.ProjectID = e.ProjectID
		case *iam.GlobalOrgSetEvent:
			wm.GlobalOrgID = e.OrgID
		case *iam.DefaultLanguageSetEvent:
			wm.DefaultLanguage = e.Language
		case *iam.SetupStepEvent:
			if e.Done {
				wm.SetUpDone = e.Step
			} else {
				wm.SetUpStarted = e.Step
			}
		}
	}
	return nil
}

func (wm *IAMWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.InstanceAddedEventType,
			iam.InstanceChangedEventType,
			iam.InstanceRemovedEventType,
			iam.InstanceDomainAddedEventType,
			iam.InstanceDomainRemovedEventType,
			iam.ProjectSetEventType,
			iam.GlobalOrgSetEventType,
			iam.DefaultLanguageSetEventType,
			iam.SetupStartedEventType,
			iam.SetupDoneEventType).
		Builder()
}

func IAMAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion)
}
