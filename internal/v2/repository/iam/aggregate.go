package iam

import (
	"context"
	iam_login "github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login/idpprovider"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

const (
	IamEventTypePrefix = eventstore.EventType("iam.")
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

func AggregateFromWriteModel(wm *eventstore.WriteModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, AggregateType, AggregateVersion),
	}
}

func AggregateFromReadModel(rm *ReadModel) *Aggregate {
	return &Aggregate{
		Aggregate: *eventstore.NewAggregate(
			rm.AggregateID,
			AggregateType,
			rm.ResourceOwner,
			AggregateVersion,
			rm.ProcessedSequence,
		),
	}
}

func (a *Aggregate) PushMemberAdded(ctx context.Context, userID string, roles ...string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberAddedEvent(ctx, userID, roles...))
	return a
}

func (a *Aggregate) PushMemberChangedFromExisting(ctx context.Context, current *MemberWriteModel, roles ...string) *Aggregate {
	e, err := MemberChangedEventFromExisting(ctx, current, roles...)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushMemberRemoved(ctx context.Context, userID string) *Aggregate {
	a.Aggregate = *a.PushEvents(NewMemberRemovedEvent(ctx, userID))
	return a
}

func (a *Aggregate) PushStepStarted(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepStartedEvent(ctx, step))
	return a
}

func (a *Aggregate) PushStepDone(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepDoneEvent(ctx, step))
	return a
}

func (a *Aggregate) PushOrgIAMPolicyAddedEvent(ctx context.Context, userLoginMustBeDomain bool) *Aggregate {
	a.Aggregate = *a.PushEvents(org_iam.NewOrgIAMPolicyAddedEventEvent(ctx, userLoginMustBeDomain))
	return a
}

func (a *Aggregate) PushOrgIAMPolicyChangedFromExisting(ctx context.Context, current *org_iam.OrgIAMPolicyWriteModel, userLoginMustBeDomain bool) *Aggregate {
	e, err := org_iam.OrgIAMPolicyChangedEventFromExisting(ctx, current, userLoginMustBeDomain)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushPasswordComplexityPolicyAddedEvent(ctx context.Context, minLength uint64, hasLowercase, hasUppercase, hasNumber, hasSymbol bool) *Aggregate {
	a.Aggregate = *a.PushEvents(password_complexity.NewPasswordComplexityPolicyAddedEvent(ctx, minLength, hasLowercase, hasUppercase, hasNumber, hasSymbol))
	return a
}

func (a *Aggregate) PushPasswordComplexityPolicyChangedFromExisting(ctx context.Context, current *password_complexity.PasswordComplexityPolicyWriteModel, minLength uint64, hasLowercase, hasUppercase, hasNumber, hasSymbol bool) *Aggregate {
	e, err := password_complexity.PasswordComplexityPolicyChangedEventFromExisting(ctx, current, minLength, hasLowercase, hasUppercase, hasNumber, hasSymbol)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
	return a
}

func (a *Aggregate) PushLoginPolicyAddedEvent(ctx context.Context, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType login.PasswordlessType) *Aggregate {
	a.Aggregate = *a.PushEvents(iam_login.NewLoginPolicyAddedEventEvent(ctx, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA, passwordlessType))
	return a
}

func (a *Aggregate) PushLoginPolicyChangedFromExisting(ctx context.Context, current *iam_login.LoginPolicyWriteModel, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType login.PasswordlessType) *Aggregate {
	e, err := iam_login.LoginPolicyChangedEventFromExisting(ctx, current, allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA, passwordlessType)
	if err != nil {
		return a
	}
	a.Aggregate = *a.PushEvents(e)
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
	providerType provider.Type,
) *Aggregate {

	a.Aggregate = *a.PushEvents(idpprovider.NewLoginPolicyIDPProviderAddedEvent(ctx, idpConfigID, providerType))
	return a
}

func (a *Aggregate) PushLoginPolicyIDPProviderRemovedEvent(
	ctx context.Context,
	idpConfigID string,
	providerType provider.Type,
) *Aggregate {

	a.Aggregate = *a.PushEvents(idpprovider.NewLoginPolicyIDPProviderRemovedEvent(ctx, idpConfigID))
	return a
}
