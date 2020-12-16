package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMWriteModel struct {
	eventstore.WriteModel

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	GlobalOrgID string
	ProjectID   string
}

func NewIAMriteModel(iamID string) *IAMWriteModel {
	return &IAMWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: iamID,
		},
	}
}

func (wm *IAMWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	//for _, event := range events {
	//	switch e := event.(type) {
	//	case *iam.LabelPolicyAddedEvent:
	//		wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
	//	case *iam.LabelPolicyChangedEvent:
	//		wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
	//	}
	//}
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
}

func (wm *IAMWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID)
}

//
//func (wm *IAMLabelPolicyWriteModel) HasChanged(primaryColor, secondaryColor string) bool {
//	if primaryColor != "" && wm.PrimaryColor != primaryColor {
//		return true
//	}
//	if secondaryColor != "" && wm.SecondaryColor != secondaryColor {
//		return true
//	}
//	return false
//}

func IAMAggregateFromWriteModel(wm *eventstore.WriteModel) *iam.Aggregate {
	return &iam.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion),
	}
}
