package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type OrgPasswordComplexityPolicyWriteModel struct {
	PasswordComplexityPolicyWriteModel
}

func NewOrgPasswordComplexityPolicyWriteModel(iamID string) *OrgPasswordComplexityPolicyWriteModel {
	return &OrgPasswordComplexityPolicyWriteModel{
		PasswordComplexityPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *OrgPasswordComplexityPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordComplexityPolicyAddedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *org.PasswordComplexityPolicyChangedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		}
	}
}

func (wm *OrgPasswordComplexityPolicyWriteModel) Reduce() error {
	return wm.PasswordComplexityPolicyWriteModel.Reduce()
}

func (wm *OrgPasswordComplexityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PasswordComplexityPolicyWriteModel.AggregateID)
}

func (wm *OrgPasswordComplexityPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) (*org.PasswordComplexityPolicyChangedEvent, bool) {

	hasChanged := false
	changedEvent := org.NewPasswordComplexityPolicyChangedEvent(ctx)
	if wm.MinLength != minLength {
		hasChanged = true
		changedEvent.MinLength = minLength
	}
	if wm.HasLowercase != hasLowercase {
		hasChanged = true
		changedEvent.HasLowercase = hasLowercase
	}
	if wm.HasUpperCase != hasUppercase {
		hasChanged = true
		changedEvent.HasUpperCase = hasUppercase
	}
	if wm.HasNumber != hasNumber {
		hasChanged = true
		changedEvent.HasNumber = hasNumber
	}
	if wm.HasSymbol != hasSymbol {
		hasChanged = true
		changedEvent.HasSymbol = hasSymbol
	}
	return changedEvent, hasChanged
}
