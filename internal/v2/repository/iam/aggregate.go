package iam

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

const (
	IAMEventTypePrefix = eventstore.EventType("iam.")
)

const (
	AggregateType    = "iam"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(
	id,
	resourceOwner string,
	previousSequence uint64,
) *Aggregate {

	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			id,
			AggregateType,
			resourceOwner,
			AggregateVersion,
			previousSequence,
		),
	}
}

//
//func (a *Aggregate) PushMemberAdded(ctx context.Context, userID string, roles ...string) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewMemberAddedEvent(ctx, userID, roles...))
//	return a
//}
//
//func (a *Aggregate) PushMemberChangedFromExisting(ctx context.Context, current *command.IAMMemberWriteModel, roles ...string) *Aggregate {
//	e, err := MemberChangedEventFromExisting(ctx, current, roles...)
//	if err != nil {
//		return a
//	}
//	a.Aggregate = *a.PushEvents(e)
//	return a
//}
//
//func (a *Aggregate) PushMemberRemoved(ctx context.Context, userID string) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewMemberRemovedEvent(ctx, userID))
//	return a
//}

func (a *Aggregate) PushStepStarted(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepStartedEvent(ctx, step))
	return a
}

func (a *Aggregate) PushStepDone(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepDoneEvent(ctx, step))
	return a
}

func (a *Aggregate) PushOrgIAMPolicyAddedEvent(ctx context.Context, userLoginMustBeDomain bool) *Aggregate {
	a.Aggregate = *a.PushEvents(org_iam.NewAddedEvent(ctx, userLoginMustBeDomain))
	return a
}

func (a *Aggregate) PushOrgIAMPolicyChangedFromExisting(ctx context.Context, current *org_iam.WriteModel, userLoginMustBeDomain bool) *Aggregate {
	e, err := org_iam.ChangedEventFromExisting(ctx, current, userLoginMustBeDomain)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushPasswordAgePolicyAddedEvent(ctx context.Context, expireWarnDays, maxAgeDays uint64) *Aggregate {
	a.Aggregate = *a.PushEvents(password_age.NewAddedEvent(ctx, expireWarnDays, maxAgeDays))
	return a
}

func (a *Aggregate) PushPasswordAgePolicyChangedFromExisting(ctx context.Context, current *password_age.WriteModel, expireWarnDays, maxAgeDays uint64) *Aggregate {
	e, err := password_age.ChangedEventFromExisting(ctx, current, expireWarnDays, maxAgeDays)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushPasswordComplexityPolicyAddedEvent(ctx context.Context, minLength uint64, hasLowercase, hasUppercase, hasNumber, hasSymbol bool) *Aggregate {
	a.Aggregate = *a.PushEvents(password_complexity.NewAddedEvent(ctx, minLength, hasLowercase, hasUppercase, hasNumber, hasSymbol))
	return a
}

func (a *Aggregate) PushPasswordComplexityPolicyChangedFromExisting(ctx context.Context, current *password_complexity.WriteModel, minLength uint64, hasLowercase, hasUppercase, hasNumber, hasSymbol bool) *Aggregate {
	e, err := password_complexity.ChangedEventFromExisting(ctx, current, minLength, hasLowercase, hasUppercase, hasNumber, hasSymbol)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushPasswordLockoutPolicyAddedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) *Aggregate {
	a.Aggregate = *a.PushEvents(password_lockout.NewAddedEvent(ctx, maxAttempts, showLockoutFailure))
	return a
}

func (a *Aggregate) PushPasswordLockoutPolicyChangedFromExisting(ctx context.Context, current *password_lockout.WriteModel, maxAttempts uint64, showLockoutFailure bool) *Aggregate {
	e, err := password_lockout.ChangedEventFromExisting(ctx, current, maxAttempts, showLockoutFailure)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

//
//func (a *Aggregate) PushLabelPolicyAddedEvent(ctx context.Context, primaryColor, secondaryColor string) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewLabelPolicyAddedEvent(ctx, primaryColor, secondaryColor))
//	return a
//}
//
//func (a *Aggregate) PushLabelPolicyChangedFromExisting(ctx context.Context, current *command.IAMLabelPolicyWriteModel, primaryColor, secondaryColor string) *Aggregate {
//	e, err := NewLabelPolicyChangedEvent(ctx, current, primaryColor, secondaryColor)
//	if err != nil {
//		return a
//	}
//	a.Aggregate = *a.PushEvents(e)
//	return a
//}
//
//func (a *Aggregate) PushLoginPolicyAddedEvent(ctx context.Context, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType domain.PasswordlessType) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewLoginPolicyAddedEvent(ctx, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA, passwordlessType))
//	return a
//}
//
//func (a *Aggregate) PushLoginPolicyChangedFromExisting(ctx context.Context, current *command.IAMLoginPolicyWriteModel, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType domain.PasswordlessType) *Aggregate {
//	e, err := NewLoginPolicyChangedEvent(ctx, current, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA, passwordlessType)
//	if err != nil {
//		return a
//	}
//	a.Aggregate = *a.PushEvents(e)
//	return a
//}
//
//func (a *Aggregate) PushLoginPolicySecondFactorAdded(ctx context.Context, mfaType domain.SecondFactorType) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewLoginPolicySecondFactorAddedEvent(ctx, mfaType))
//	return a
//}
//
//func (a *Aggregate) PushLoginPolicySecondFactorRemoved(ctx context.Context, mfaType domain.SecondFactorType) *Aggregate {
//	a.Aggregate = *a.PushEvents(NewLoginPolicySecondFactorRemovedEvent(ctx, mfaType))
//	return a
//}

func (a *Aggregate) PushLoginPolicyMultiFactorAdded(ctx context.Context, mfaType domain.MultiFactorType) *Aggregate {
	a.Aggregate = *a.PushEvents(NewLoginPolicyMultiFactorAddedEvent(ctx, mfaType))
	return a
}

func (a *Aggregate) PushLoginPolicyMultiFactorRemoved(ctx context.Context, mfaType domain.MultiFactorType) *Aggregate {
	a.Aggregate = *a.PushEvents(NewLoginPolicyMultiFactorRemovedEvent(ctx, mfaType))
	return a
}

func (a *Aggregate) PushIDPConfigAdded(
	ctx context.Context,
	configID,
	name string,
	configType idp.ConfigType,
	stylingType idp.StylingType,
) *Aggregate {

	a.Aggregate = *a.PushEvents(NewIDPConfigAddedEvent(ctx, configID, name, configType, stylingType))
	return a
}

func (a *Aggregate) PushIDPConfigChanged(
	ctx context.Context,
	current *IDPConfigWriteModel,
	configID,
	name string,
	configType idp.ConfigType,
	stylingType idp.StylingType,
) *Aggregate {

	event, err := NewIDPConfigChangedEvent(ctx, current, configID, name, configType, stylingType)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(event)
	return a
}

func (a *Aggregate) PushIDPConfigDeactivated(ctx context.Context, configID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewIDPConfigDeactivatedEvent(ctx, configID))
	return a
}

func (a *Aggregate) PushIDPConfigReactivated(ctx context.Context, configID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewIDPConfigReactivatedEvent(ctx, configID))
	return a
}

func (a *Aggregate) PushIDPConfigRemoved(ctx context.Context, configID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewIDPConfigRemovedEvent(ctx, configID))
	return a
}

func (a *Aggregate) PushIDPOIDCConfigAdded(
	ctx context.Context,
	clientID,
	idpConfigID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping oidc.MappingField,
	scopes ...string,
) *Aggregate {

	a.Aggregate = *a.PushEvents(NewIDPOIDCConfigAddedEvent(ctx, clientID, idpConfigID, issuer, clientSecret, idpDisplayNameMapping, userNameMapping, scopes...))
	return a
}

func (a *Aggregate) PushIDPOIDCConfigChanged(
	ctx context.Context,
	current *IDPOIDCConfigWriteModel,
	clientID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping oidc.MappingField,
	scopes ...string,
) *Aggregate {

	event, err := NewIDPOIDCConfigChangedEvent(ctx, current, clientID, issuer, clientSecret, idpDisplayNameMapping, userNameMapping, scopes...)
	if err != nil {
		return a
	}

	a.Aggregate = *a.PushEvents(event)
	return a
}

func (a *Aggregate) PushLoginPolicyIDPProviderAddedEvent(
	ctx context.Context,
	idpConfigID string,
	providerType domain.IdentityProviderType,
) *Aggregate {

	a.Aggregate = *a.PushEvents(NewIAMIdentityProviderAddedEvent(ctx, idpConfigID, providerType))
	return a
}

func (a *Aggregate) PushLoginPolicyIDPProviderRemovedEvent(
	ctx context.Context,
	idpConfigID string,
	providerType domain.IdentityProviderType,
) *Aggregate {

	a.Aggregate = *a.PushEvents(NewIAMIdentityProviderRemovedEvent(ctx, idpConfigID))
	return a
}
