package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
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
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IAMPasswordComplexityPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) (*iam.PasswordComplexityPolicyChangedEvent, bool) {

	hasChanged := false
	changedEvent := iam.NewPasswordComplexityPolicyChangedEvent(ctx)
	if wm.MinLength != minLength {
		hasChanged = true
		changedEvent.MinLength = &minLength
	}
	if wm.HasLowercase != hasLowercase {
		hasChanged = true
		changedEvent.HasLowercase = &hasLowercase
	}
	if wm.HasUpperCase != hasUppercase {
		hasChanged = true
		changedEvent.HasUpperCase = &hasUppercase
	}
	if wm.HasNumber != hasNumber {
		hasChanged = true
		changedEvent.HasNumber = &hasNumber
	}
	if wm.HasSymbol != hasSymbol {
		hasChanged = true
		changedEvent.HasSymbol = &hasSymbol
	}
	return changedEvent, hasChanged
}
