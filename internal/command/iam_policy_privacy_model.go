package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMPrivacyPolicyWriteModel struct {
	PrivacyPolicyWriteModel
}

func NewIAMPrivacyPolicyWriteModel() *IAMPrivacyPolicyWriteModel {
	return &IAMPrivacyPolicyWriteModel{
		PrivacyPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMPrivacyPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PrivacyPolicyAddedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyAddedEvent)
		case *iam.PrivacyPolicyChangedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyChangedEvent)
		}
	}
}

func (wm *IAMPrivacyPolicyWriteModel) Reduce() error {
	return wm.PrivacyPolicyWriteModel.Reduce()
}

func (wm *IAMPrivacyPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PrivacyPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.PrivacyPolicyAddedEventType,
			iam.PrivacyPolicyChangedEventType)
}

func (wm *IAMPrivacyPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tosLink,
	privacyLink string,
) (*iam.PrivacyPolicyChangedEvent, bool) {

	changes := make([]policy.PrivacyPolicyChanges, 0)
	if wm.TOSLink != tosLink {
		changes = append(changes, policy.ChangeTOSLink(tosLink))
	}
	if wm.PrivacyLink != privacyLink {
		changes = append(changes, policy.ChangePrivacyLink(privacyLink))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewPrivacyPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
