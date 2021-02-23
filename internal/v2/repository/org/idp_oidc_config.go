package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

const (
	IDPOIDCConfigAddedEventType   eventstore.EventType = "org.idp." + idpconfig.OIDCConfigAddedEventType
	IDPOIDCConfigChangedEventType eventstore.EventType = "org.idp." + idpconfig.ConfigChangedEventType
)

type IDPOIDCConfigAddedEvent struct {
	idpconfig.OIDCConfigAddedEvent
}

func NewIDPOIDCConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID,
	idpConfigID,
	issuer string,
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
