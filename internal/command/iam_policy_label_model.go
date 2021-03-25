package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMLabelPolicyWriteModel struct {
	LabelPolicyWriteModel
}

func NewIAMLabelPolicyWriteModel() *IAMLabelPolicyWriteModel {
	return &IAMLabelPolicyWriteModel{
		LabelPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMLabelPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LabelPolicyAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *iam.LabelPolicyChangedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
		}
	}
}

func (wm *IAMLabelPolicyWriteModel) Reduce() error {
	return wm.LabelPolicyWriteModel.Reduce()
}

func (wm *IAMLabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.LabelPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.LabelPolicyAddedEventType,
			iam.LabelPolicyChangedEventType)
}

func (wm *IAMLabelPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	primaryColor,
	secondaryColor string,
	hideLoginNameSuffix bool,
) (*iam.LabelPolicyChangedEvent, bool) {
	changes := make([]policy.LabelPolicyChanges, 0)
	if wm.PrimaryColor != primaryColor {
		changes = append(changes, policy.ChangePrimaryColor(primaryColor))
	}
	if wm.SecondaryColor != secondaryColor {
		changes = append(changes, policy.ChangeSecondaryColor(secondaryColor))
	}
	if wm.SecondaryColor != secondaryColor {
		changes = append(changes, policy.ChangeSecondaryColor(secondaryColor))
	}
	if wm.HideLoginNameSuffix != hideLoginNameSuffix {
		changes = append(changes, policy.ChangeHideLoginNameSuffix(hideLoginNameSuffix))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewLabelPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
