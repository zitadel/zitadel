package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type InstanceLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewInstanceLoginPolicyWriteModel(ctx context.Context) *InstanceLoginPolicyWriteModel {
	return &InstanceLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceLoginPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *instance.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		}
	}
}

func (wm *InstanceLoginPolicyWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *InstanceLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *InstanceLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		EventTypes(
			instance.LoginPolicyAddedEventType,
			instance.LoginPolicyChangedEventType).
		Builder()
}

func (wm *InstanceLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA,
	forceMFALocalOnly,
	hidePasswordReset,
	ignoreUnknownUsernames,
	allowDomainDiscovery,
	disableLoginWithEmail,
	disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	defaultRedirectURI string,
	passwordCheckLifetime,
	externalLoginCheckLifetime,
	mfaInitSkipLifetime,
	secondFactorCheckLifetime,
	multiFactorCheckLifetime time.Duration,
) (*instance.LoginPolicyChangedEvent, bool) {

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
	if wm.ForceMFALocalOnly != forceMFALocalOnly {
		changes = append(changes, policy.ChangeForceMFALocalOnly(forceMFALocalOnly))
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		changes = append(changes, policy.ChangePasswordlessType(passwordlessType))
	}
	if wm.HidePasswordReset != hidePasswordReset {
		changes = append(changes, policy.ChangeHidePasswordReset(hidePasswordReset))
	}
	if wm.IgnoreUnknownUsernames != ignoreUnknownUsernames {
		changes = append(changes, policy.ChangeIgnoreUnknownUsernames(ignoreUnknownUsernames))
	}
	if wm.AllowDomainDiscovery != allowDomainDiscovery {
		changes = append(changes, policy.ChangeAllowDomainDiscovery(allowDomainDiscovery))
	}
	if wm.DefaultRedirectURI != defaultRedirectURI {
		changes = append(changes, policy.ChangeDefaultRedirectURI(defaultRedirectURI))
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
	if wm.DisableLoginWithEmail != disableLoginWithEmail {
		changes = append(changes, policy.ChangeDisableLoginWithEmail(disableLoginWithEmail))
	}
	if wm.DisableLoginWithPhone != disableLoginWithPhone {
		changes = append(changes, policy.ChangeDisableLoginWithPhone(disableLoginWithPhone))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewLoginPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
