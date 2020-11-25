package iam

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

const (
	IDPOIDCConfigAddedEventType   eventstore.EventType = "iam.idp.oidc.config.added"
	IDPOIDCConfigChangedEventType eventstore.EventType = "iam.idp.oidc.config.changed"
)

type IDPOIDCConfigWriteModel struct {
	oidc.ConfigWriteModel
}

func (rm *IDPOIDCConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPOIDCConfigAddedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPOIDCConfigChangedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigChangedEvent)
		case *oidc.ConfigAddedEvent,
			*oidc.ConfigChangedEvent:

			rm.ConfigWriteModel.AppendEvents(e)
		}
	}
}

type IDPOIDCConfigAddedEvent struct {
	oidc.ConfigAddedEvent
}

func NewIDPOIDCConfigAddedEvent(
	ctx context.Context,
	clientID,
	idpConfigID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping oidc.MappingField,
	scopes ...string,
) *IDPOIDCConfigAddedEvent {

	return &IDPOIDCConfigAddedEvent{
		ConfigAddedEvent: *oidc.NewConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
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

type IDPOIDCConfigChangedEvent struct {
	oidc.ConfigChangedEvent
}

func NewIDPOIDCConfigChangedEvent(
	ctx context.Context,
	current *IDPOIDCConfigWriteModel,
	clientID,
	idpConfigID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping oidc.MappingField,
	scopes ...string,
) (*IDPOIDCConfigChangedEvent, error) {

	event, err := oidc.NewConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			IDPOIDCConfigAddedEventType,
		),
		&current.ConfigWriteModel,
		clientID,
		issuer,
		clientSecret,
		idpDisplayNameMapping,
		userNameMapping,
		scopes...,
	)

	if err != nil {
		return nil, err
	}

	return &IDPOIDCConfigChangedEvent{
		ConfigChangedEvent: *event,
	}, nil
}
