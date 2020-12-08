package idpprovider

import (
	"context"
	"github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type LoginPolicyIDPProviderAddedEvent struct {
	idpprovider.IDPProviderAddedEvent
}

func NewLoginPolicyIDPProviderAddedEvent(
	ctx context.Context,
	idpConfigID string,
	idpProviderType provider.Type,
) *LoginPolicyIDPProviderAddedEvent {

	return &LoginPolicyIDPProviderAddedEvent{
		IDPProviderAddedEvent: *idpprovider.NewIDPProviderAddedEvent(
			eventstore.NewBaseEventForPush(ctx, login.LoginPolicyIDPProviderAddedEventType),
			idpConfigID,
			idpProviderType),
	}
}

func IDPProviderAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpprovider.IDPProviderAddedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyIDPProviderAddedEvent{
		IDPProviderAddedEvent: *e.(*idpprovider.IDPProviderAddedEvent),
	}, nil
}

type LoginPolicyIDPProviderRemovedEvent struct {
	idpprovider.IDPProviderRemovedEvent
}

func NewLoginPolicyIDPProviderRemovedEvent(
	ctx context.Context,
	idpConfigID string,
) *LoginPolicyIDPProviderRemovedEvent {

	return &LoginPolicyIDPProviderRemovedEvent{
		IDPProviderRemovedEvent: *idpprovider.NewIDPProviderRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, login.LoginPolicyIDPProviderRemovedEventType),
			idpConfigID),
	}
}

func IDPProviderRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpprovider.IDPProviderRemovedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyIDPProviderRemovedEvent{
		IDPProviderRemovedEvent: *e.(*idpprovider.IDPProviderRemovedEvent),
	}, nil
}
