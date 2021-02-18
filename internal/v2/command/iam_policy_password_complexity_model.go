package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMPasswordComplexityPolicyWriteModel struct {
	PasswordComplexityPolicyWriteModel
}

func NewIAMPasswordComplexityPolicyWriteModel() *IAMPasswordComplexityPolicyWriteModel {
	return &IAMPasswordComplexityPolicyWriteModel{
		PasswordComplexityPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMPasswordComplexityPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordComplexityPolicyAddedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *iam.PasswordComplexityPolicyChangedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		}
	}
}

func (wm *IAMPasswordComplexityPolicyWriteModel) Reduce() error {
	return wm.PasswordComplexityPolicyWriteModel.Reduce()
}

func (wm *IAMPasswordComplexityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PasswordComplexityPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.PasswordComplexityPolicyAddedEventType,
			iam.PasswordComplexityPolicyChangedEventType)
}

func (wm *IAMPasswordComplexityPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) (*iam.PasswordComplexityPolicyChangedEvent, bool) {

	changes := make([]policy.PasswordComplexityPolicyChanges, 0)
	if wm.MinLength != minLength {
		changes = append(changes, policy.ChangeMinLength(minLength))
	}
	if wm.HasLowercase != hasLowercase {
		changes = append(changes, policy.ChangeHasLowercase(hasLowercase))
	}
	if wm.HasUppercase != hasUppercase {
		changes = append(changes, policy.ChangeHasUppercase(hasUppercase))
	}
	if wm.HasNumber != hasNumber {
		changes = append(changes, policy.ChangeHasNumber(hasNumber))
	}
	if wm.HasSymbol != hasSymbol {
		changes = append(changes, policy.ChangeHasSymbol(hasSymbol))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewPasswordComplexityPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
