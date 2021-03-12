package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/features"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMFeaturesWriteModel struct {
	FeaturesWriteModel
}

func NewIAMFeaturesWriteModel() *IAMFeaturesWriteModel {
	return &IAMFeaturesWriteModel{
		FeaturesWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMFeaturesWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.FeaturesSetEvent:
			wm.FeaturesWriteModel.AppendEvents(&e.FeaturesSetEvent)
		}
	}
}

func (wm *IAMFeaturesWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *IAMFeaturesWriteModel) Reduce() error {
	return wm.FeaturesWriteModel.Reduce()
}

func (wm *IAMFeaturesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.FeaturesWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(iam.FeaturesSetEventType)
}

func (wm *IAMFeaturesWriteModel) NewSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tierName,
	tierDescription string,
	tierStatus domain.TierStatus,
	tierStatusDescription string,
	loginPolicyFactors,
	loginPolicyIDP,
	loginPolicyPasswordless,
	loginPolicyRegistration,
	loginPolicyUsernameLogin bool,
) (*iam.FeaturesSetEvent, bool) {

	changes := make([]features.FeaturesChanges, 0)

	if tierName != "" && wm.TierName != tierName {
		changes = append(changes, features.ChangeTierName(tierName))
	}
	if tierDescription != "" && wm.TierDescription != tierDescription {
		changes = append(changes, features.ChangeTierDescription(tierDescription))
	}
	if wm.TierStatus != tierStatus {
		changes = append(changes, features.ChangeTierStatus(tierStatus))
	}
	if tierStatusDescription != "" && wm.TierStatusDescription != tierStatusDescription {
		changes = append(changes, features.ChangeTierStatusDescription(tierStatusDescription))
	}
	if wm.LoginPolicyFactors != loginPolicyFactors {
		changes = append(changes, features.ChangeLoginPolicyFactors(loginPolicyFactors))
	}
	if wm.LoginPolicyIDP != loginPolicyIDP {
		changes = append(changes, features.ChangeLoginPolicyIDP(loginPolicyIDP))
	}
	if wm.LoginPolicyPasswordless != loginPolicyPasswordless {
		changes = append(changes, features.ChangeLoginPolicyPasswordless(loginPolicyPasswordless))
	}
	if wm.LoginPolicyRegistration != loginPolicyRegistration {
		changes = append(changes, features.ChangeLoginPolicyRegistration(loginPolicyRegistration))
	}
	if wm.LoginPolicyUsernameLogin != loginPolicyUsernameLogin {
		changes = append(changes, features.ChangeLoginPolicyUsernameLogin(loginPolicyUsernameLogin))
	}

	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewFeaturesSetEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
