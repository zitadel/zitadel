package idpprovider

import (
	"context"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

var (
	iamEventPrefix = eventstore.EventType("iam.")

	LoginPolicyIDPProviderAddedEventType   = iamEventPrefix + idpprovider.LoginPolicyIDPProviderAddedType
	LoginPolicyIDPProviderRemovedEventType = iamEventPrefix + idpprovider.LoginPolicyIDPProviderRemovedType
)

type AddedEvent struct {
	idpprovider.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	idpConfigID string,
	idpProviderType idpprovider.Type,
) *AddedEvent {

	return &AddedEvent{
		AddedEvent: *idpprovider.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyIDPProviderAddedEventType),
			idpConfigID,
			idpProviderType),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpprovider.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{
		AddedEvent: *e.(*idpprovider.AddedEvent),
	}, nil
}

type RemovedEvent struct {
	idpprovider.RemovedEvent
}

func NewRemovedEvent(
	ctx context.Context,
	idpConfigID string,
) *RemovedEvent {
	return &RemovedEvent{
		RemovedEvent: *idpprovider.NewRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyIDPProviderRemovedEventType),
			idpConfigID),
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpprovider.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &RemovedEvent{
		RemovedEvent: *e.(*idpprovider.RemovedEvent),
	}, nil
}
