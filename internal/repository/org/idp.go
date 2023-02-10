package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

const (
	GoogleIDPAddedEventType   eventstore.EventType = "org.idp.google.added"
	GoogleIDPChangedEventType eventstore.EventType = "org.idp.google.changed"
	OAuthIDPAddedEventType    eventstore.EventType = "org.idp.oauth.added"
	OAuthIDPChangedEventType  eventstore.EventType = "org.idp.oauth.changed"
	GitHubIDPAddedEventType   eventstore.EventType = "org.idp.github.added"
	GitHubIDPChangedEventType eventstore.EventType = "org.idp.github.changed"
)

type GoogleIDPAddedEvent struct {
	idp.GoogleIDPAddedEvent
}

func NewGoogleIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	options idp.Options,
) *GoogleIDPAddedEvent {

	return &GoogleIDPAddedEvent{
		GoogleIDPAddedEvent: *idp.NewGoogleIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GoogleIDPAddedEventType,
			),
			id,
			clientID,
			clientSecret,
			options,
		),
	}
}

func GoogleIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GoogleIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GoogleIDPAddedEvent{GoogleIDPAddedEvent: *e.(*idp.GoogleIDPAddedEvent)}, nil
}

type GoogleIDPChangedEvent struct {
	idp.GoogleIDPChangedEvent
}

func NewGoogleIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GoogleIDPChanges,
) (*GoogleIDPChangedEvent, error) {

	changedEvent, err := idp.NewGoogleIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GoogleIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GoogleIDPChangedEvent{GoogleIDPChangedEvent: *changedEvent}, nil
}

func GoogleIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GoogleIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GoogleIDPChangedEvent{GoogleIDPChangedEvent: *e.(*idp.GoogleIDPChangedEvent)}, nil
}

type OAuthIDPAddedEvent struct {
	idp.OAuthIDPAddedEvent
}

func NewOAuthIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) *OAuthIDPAddedEvent {

	return &OAuthIDPAddedEvent{
		OAuthIDPAddedEvent: *idp.NewOAuthIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OAuthIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			authorizationEndpoint,
			tokenEndpoint,
			userEndpoint,
			scopes,
			options,
		),
	}
}

func OAuthIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OAuthIDPAddedEvent{OAuthIDPAddedEvent: *e.(*idp.OAuthIDPAddedEvent)}, nil
}

type OAuthIDPChangedEvent struct {
	idp.OAuthIDPChangedEvent
}

func NewOAuthIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName string,
	changes []idp.OAuthIDPChanges,
) (*OAuthIDPChangedEvent, error) {

	changedEvent, err := idp.NewOAuthIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OAuthIDPChangedEventType,
		),
		id,
		oldName,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OAuthIDPChangedEvent{OAuthIDPChangedEvent: *changedEvent}, nil
}

func OAuthIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OAuthIDPChangedEvent{OAuthIDPChangedEvent: *e.(*idp.OAuthIDPChangedEvent)}, nil
}

type GitHubIDPAddedEvent struct {
	idp.GitHubIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	options idp.Options,
) *GitHubIDPAddedEvent {

	return &GitHubIDPAddedEvent{
		GitHubIDPAddedEvent: *idp.NewGitHubIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitHubIDPAddedEventType,
			),
			id,
			clientID,
			clientSecret,
			options,
		),
	}
}

func GitHubIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GitHubIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPAddedEvent{GitHubIDPAddedEvent: *e.(*idp.GitHubIDPAddedEvent)}, nil
}

type GitHubIDPChangedEvent struct {
	idp.GitHubIDPChangedEvent
}

func NewGitHubIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.OAuthIDPChanges,
) (*GitHubIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitHubIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitHubIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitHubIDPChangedEvent{GitHubIDPChangedEvent: *changedEvent}, nil
}

func GitHubIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GitHubIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPChangedEvent{GitHubIDPChangedEvent: *e.(*idp.GitHubIDPChangedEvent)}, nil
}
