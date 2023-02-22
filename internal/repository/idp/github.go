package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type GitHubIDPAddedEvent struct {
	OAuthIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GitHubIDPAddedEvent {
	return &GitHubIDPAddedEvent{
		OAuthIDPAddedEvent: OAuthIDPAddedEvent{
			BaseEvent:    *base,
			ID:           id,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Options:      options,
		},
	}
}

func (e *GitHubIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitHubIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPAddedEvent{OAuthIDPAddedEvent: *e.(*OAuthIDPAddedEvent)}, nil
}

type GitHubIDPChangedEvent struct {
	OAuthIDPChangedEvent
}

func NewGitHubIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*GitHubIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubIDPChangedEvent{
		OAuthIDPChangedEvent: OAuthIDPChangedEvent{
			BaseEvent: *base,
			ID:        id,
		},
	}
	for _, change := range changes {
		change(&changedEvent.OAuthIDPChangedEvent)
	}
	return changedEvent, nil
}

func (e *GitHubIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitHubIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPChangedEvent{OAuthIDPChangedEvent: *e.(*OAuthIDPChangedEvent)}, nil
}

type GitHubEnterpriseIDPAddedEvent struct {
	OAuthIDPAddedEvent
}

func NewGitHubEnterpriseIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options Options,
) *GitHubEnterpriseIDPAddedEvent {
	return &GitHubEnterpriseIDPAddedEvent{
		OAuthIDPAddedEvent: *NewOAuthIDPAddedEvent(
			base,
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

func (e *GitHubEnterpriseIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitHubEnterpriseIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitHubEnterpriseIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPAddedEvent{OAuthIDPAddedEvent: *e.(*OAuthIDPAddedEvent)}, nil
}

type GitHubEnterpriseIDPChangedEvent struct {
	OAuthIDPChangedEvent
}

func NewGitHubEnterpriseIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*GitHubEnterpriseIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-JHKs9", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubEnterpriseIDPChangedEvent{
		OAuthIDPChangedEvent: OAuthIDPChangedEvent{
			BaseEvent: *base,
			ID:        id,
		},
	}
	for _, change := range changes {
		change(&changedEvent.OAuthIDPChangedEvent)
	}
	return changedEvent, nil
}

func (e *GitHubEnterpriseIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitHubEnterpriseIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitHubEnterpriseIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPChangedEvent{OAuthIDPChangedEvent: *e.(*OAuthIDPChangedEvent)}, nil
}
