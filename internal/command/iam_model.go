package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMWriteModel struct {
	eventstore.WriteModel

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	GlobalOrgID string
	ProjectID   string
}

func NewIAMWriteModel() *IAMWriteModel {
	return &IAMWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.ProjectSetEvent:
			wm.ProjectID = e.ProjectID
		case *iam.GlobalOrgSetEvent:
			wm.GlobalOrgID = e.OrgID
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
			iam.ProjectSetEventType,
			iam.GlobalOrgSetEventType,
			iam.SetupStartedEventType,
			iam.SetupDoneEventType).
		Builder()
}

func IAMAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion)
}
