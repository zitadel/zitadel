package iam

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

const (
	iamEventTypePrefix = eventstore.EventType("iam.")
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

func (a *Aggregate) PushStepStarted(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepStartedEvent(ctx, step))
	return a
}

func (a *Aggregate) PushStepDone(ctx context.Context, step Step) *Aggregate {
	a.Aggregate = *a.PushEvents(NewSetupStepDoneEvent(ctx, step))
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
