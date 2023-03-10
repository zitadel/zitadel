package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type GitLabIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         string              `json:"name,omitempty"`
	ClientID     string              `json:"client_id"`
	ClientSecret *crypto.CryptoValue `json:"client_secret"`
	Scopes       []string            `json:"scopes,omitempty"`
	Options
}

func NewGitLabIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GitLabIDPAddedEvent {
	return &GitLabIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *GitLabIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitLabIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitLabIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GitLabIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-KLewio", "unable to unmarshal event")
	}

	return e, nil
}

type GitLabIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         *string             `json:"name,omitempty"`
	ClientID     *string             `json:"client_id,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"client_secret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewGitLabIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GitLabIDPChanges,
) (*GitLabIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-K2gje", "Errors.NoChangesFound")
	}
	changedEvent := &GitLabIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type GitLabIDPChanges func(*GitLabIDPChangedEvent)

func ChangeGitLabName(name string) func(*GitLabIDPChangedEvent) {
	return func(e *GitLabIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeGitLabClientID(clientID string) func(*GitLabIDPChangedEvent) {
	return func(e *GitLabIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeGitLabClientSecret(clientSecret *crypto.CryptoValue) func(*GitLabIDPChangedEvent) {
	return func(e *GitLabIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeGitLabScopes(scopes []string) func(*GitLabIDPChangedEvent) {
	return func(e *GitLabIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeGitLabOptions(options OptionChanges) func(*GitLabIDPChangedEvent) {
	return func(e *GitLabIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *GitLabIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitLabIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitLabIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GitLabIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sfhjk", "unable to unmarshal event")
	}

	return e, nil
}

type GitLabSelfHostedIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Issuer       string              `json:"issuer"`
	ClientID     string              `json:"client_id"`
	ClientSecret *crypto.CryptoValue `json:"client_secret"`
	Scopes       []string            `json:"scopes,omitempty"`
	Options
}

func NewGitLabSelfHostedIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GitLabSelfHostedIDPAddedEvent {
	return &GitLabSelfHostedIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		Issuer:       issuer,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *GitLabSelfHostedIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitLabSelfHostedIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitLabSelfHostedIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GitLabSelfHostedIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-S1efv", "unable to unmarshal event")
	}

	return e, nil
}

type GitLabSelfHostedIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         *string             `json:"name,omitempty"`
	Issuer       *string             `json:"issuer,omitempty"`
	ClientID     *string             `json:"client_id,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"client_secret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewGitLabSelfHostedIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GitLabSelfHostedIDPChanges,
) (*GitLabSelfHostedIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-Dghj6", "Errors.NoChangesFound")
	}
	changedEvent := &GitLabSelfHostedIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type GitLabSelfHostedIDPChanges func(*GitLabSelfHostedIDPChangedEvent)

func ChangeGitLabSelfHostedName(name string) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeGitLabSelfHostedIssuer(issuer string) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeGitLabSelfHostedClientID(clientID string) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeGitLabSelfHostedClientSecret(clientSecret *crypto.CryptoValue) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeGitLabSelfHostedScopes(scopes []string) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeGitLabSelfHostedOptions(options OptionChanges) func(*GitLabSelfHostedIDPChangedEvent) {
	return func(e *GitLabSelfHostedIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *GitLabSelfHostedIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitLabSelfHostedIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func GitLabSelfHostedIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GitLabSelfHostedIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-SFrhj", "unable to unmarshal event")
	}

	return e, nil
}
