package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

const (
	IDPOIDCConfigAddedEventType   eventstore.EventType = "iam.idp." + idpconfig.OIDCConfigAddedEventType
	IDPOIDCConfigChangedEventType eventstore.EventType = "iam.idp." + idpconfig.ConfigChangedEventType
)

type IDPOIDCConfigAddedEvent struct {
	idpconfig.OIDCConfigAddedEvent
}

func NewIDPOIDCConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID,
	idpConfigID,
	issuer,
	authorizationEndpoint,
	tokenEndpoint string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) *IDPOIDCConfigAddedEvent {

	return &IDPOIDCConfigAddedEvent{
		OIDCConfigAddedEvent: *idpconfig.NewOIDCConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPOIDCConfigAddedEventType,
			),
			clientID,
			idpConfigID,
			issuer,
			authorizationEndpoint,
			tokenEndpoint,
			clientSecret,
			idpDisplayNameMapping,
			userNameMapping,
			scopes...,
		),
	}
}

func IDPOIDCConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.OIDCConfigAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPOIDCConfigAddedEvent{OIDCConfigAddedEvent: *e.(*idpconfig.OIDCConfigAddedEvent)}, nil
}

type IDPOIDCConfigChangedEvent struct {
	idpconfig.OIDCConfigChangedEvent
}

func NewIDPOIDCConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
	changes []idpconfig.OIDCConfigChanges,
) (*IDPOIDCConfigChangedEvent, error) {
	changeEvent, err := idpconfig.NewOIDCConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IDPOIDCConfigChangedEventType),
		idpConfigID,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &IDPOIDCConfigChangedEvent{OIDCConfigChangedEvent: *changeEvent}, nil
}

func IDPOIDCConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.OIDCConfigChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPOIDCConfigChangedEvent{OIDCConfigChangedEvent: *e.(*idpconfig.OIDCConfigChangedEvent)}, nil
}
