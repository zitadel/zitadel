package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgPasswordComplexityPolicyWriteModel struct {
	PasswordComplexityPolicyWriteModel
}

func NewOrgPasswordComplexityPolicyWriteModel(orgID string) *OrgPasswordComplexityPolicyWriteModel {
	return &OrgPasswordComplexityPolicyWriteModel{
		PasswordComplexityPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
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
		case *org.PasswordComplexityPolicyRemovedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyRemovedEvent)
		}
	}
}

func (wm *OrgPasswordComplexityPolicyWriteModel) Reduce() error {
	return wm.PasswordComplexityPolicyWriteModel.Reduce()
}

func (wm *OrgPasswordComplexityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PasswordComplexityPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgPasswordComplexityPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) (*org.PasswordComplexityPolicyChangedEvent, bool) {

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
	changedEvent, err := org.NewPasswordComplexityPolicyChangedEvent(ctx, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
