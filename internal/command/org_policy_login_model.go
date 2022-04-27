package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type OrgLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewOrgLoginPolicyWriteModel(orgID string) *OrgLoginPolicyWriteModel {
	return &OrgLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLoginPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *org.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *org.LoginPolicyRemovedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyRemovedEvent)
		}
	}
}

func (wm *OrgLoginPolicyWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *OrgLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *OrgLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		EventTypes(
			org.LoginPolicyAddedEventType,
			org.LoginPolicyChangedEventType,
			org.LoginPolicyRemovedEventType).
		Builder()
}

func (wm *OrgLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA,
	hidePasswordReset bool,
	passwordlessType domain.PasswordlessType,
	passwordCheckLifetime,
	externalLoginCheckLifetime,
	mfaInitSkipLifetime,
	secondFactorCheckLifetime,
	multiFactorCheckLifetime time.Duration,
) (*org.LoginPolicyChangedEvent, bool) {

	changes := make([]policy.LoginPolicyChanges, 0)
	if wm.AllowUserNamePassword != allowUsernamePassword {
		changes = append(changes, policy.ChangeAllowUserNamePassword(allowUsernamePassword))
	}
	if wm.AllowRegister != allowRegister {
		changes = append(changes, policy.ChangeAllowRegister(allowRegister))
	}
	if wm.AllowExternalIDP != allowExternalIDP {
		changes = append(changes, policy.ChangeAllowExternalIDP(allowExternalIDP))
	}
	if wm.ForceMFA != forceMFA {
		changes = append(changes, policy.ChangeForceMFA(forceMFA))
	}
	if wm.HidePasswordReset != hidePasswordReset {
		changes = append(changes, policy.ChangeHidePasswordReset(hidePasswordReset))
	}
	if wm.PasswordCheckLifetime != passwordCheckLifetime {
		changes = append(changes, policy.ChangePasswordCheckLifetime(passwordCheckLifetime))
	}
	if wm.ExternalLoginCheckLifetime != externalLoginCheckLifetime {
		changes = append(changes, policy.ChangeExternalLoginCheckLifetime(externalLoginCheckLifetime))
	}
	if wm.MFAInitSkipLifetime != mfaInitSkipLifetime {
		changes = append(changes, policy.ChangeMFAInitSkipLifetime(mfaInitSkipLifetime))
	}
	if wm.SecondFactorCheckLifetime != secondFactorCheckLifetime {
		changes = append(changes, policy.ChangeSecondFactorCheckLifetime(secondFactorCheckLifetime))
	}
	if wm.MultiFactorCheckLifetime != multiFactorCheckLifetime {
		changes = append(changes, policy.ChangeMultiFactorCheckLifetime(multiFactorCheckLifetime))
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		changes = append(changes, policy.ChangePasswordlessType(passwordlessType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewLoginPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
