package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyIDPProviderAddedEventType          = iamEventTypePrefix + policy.LoginPolicyIDPProviderAddedType
	LoginPolicyIDPProviderRemovedEventType        = iamEventTypePrefix + policy.LoginPolicyIDPProviderRemovedType
	LoginPolicyIDPProviderCascadeRemovedEventType = iamEventTypePrefix + policy.LoginPolicyIDPProviderCascadeRemovedType
)

type IdentityProviderAddedEvent struct {
	policy.IdentityProviderAddedEvent
}

func NewIdentityProviderAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
	idpProviderType domain.IdentityProviderType,
) *IdentityProviderAddedEvent {

	return &IdentityProviderAddedEvent{
		IdentityProviderAddedEvent: *policy.NewIdentityProviderAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicyIDPProviderAddedEventType),
			idpConfigID,
			idpProviderType),
	}
}

func IdentityProviderAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.IdentityProviderAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IdentityProviderAddedEvent{
		IdentityProviderAddedEvent: *e.(*policy.IdentityProviderAddedEvent),
	}, nil
}

type IdentityProviderRemovedEvent struct {
	policy.IdentityProviderRemovedEvent
}

func NewIdentityProviderRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
) *IdentityProviderRemovedEvent {
	return &IdentityProviderRemovedEvent{
		IdentityProviderRemovedEvent: *policy.NewIdentityProviderRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicyIDPProviderRemovedEventType),
			idpConfigID),
	}
}

func IdentityProviderRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.IdentityProviderRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IdentityProviderRemovedEvent{
		IdentityProviderRemovedEvent: *e.(*policy.IdentityProviderRemovedEvent),
	}, nil
}

type IdentityProviderCascadeRemovedEvent struct {
	policy.IdentityProviderCascadeRemovedEvent
}

func NewIdentityProviderCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
) *IdentityProviderCascadeRemovedEvent {
	return &IdentityProviderCascadeRemovedEvent{
		IdentityProviderCascadeRemovedEvent: *policy.NewIdentityProviderCascadeRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, LoginPolicyIDPProviderCascadeRemovedEventType),
			idpConfigID),
	}
}

func IdentityProviderCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.IdentityProviderCascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IdentityProviderCascadeRemovedEvent{
		IdentityProviderCascadeRemovedEvent: *e.(*policy.IdentityProviderCascadeRemovedEvent),
	}, nil
}
