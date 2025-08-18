package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GitHubIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         string              `json:"name,omitempty"`
	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	Options
}

func NewGitHubIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GitHubIDPAddedEvent {
	return &GitHubIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *GitHubIDPAddedEvent) Payload() any {
	return e
}

func (e *GitHubIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GitHubIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GitHubIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-sdfs3", "unable to unmarshal event")
	}

	return e, nil
}

type GitHubIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         *string             `json:"name,omitempty"`
	ClientID     *string             `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewGitHubIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GitHubIDPChanges,
) (*GitHubIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type GitHubIDPChanges func(*GitHubIDPChangedEvent)

func ChangeGitHubName(name string) func(*GitHubIDPChangedEvent) {
	return func(e *GitHubIDPChangedEvent) {
		e.Name = &name
	}
}
func ChangeGitHubClientID(clientID string) func(*GitHubIDPChangedEvent) {
	return func(e *GitHubIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeGitHubClientSecret(clientSecret *crypto.CryptoValue) func(*GitHubIDPChangedEvent) {
	return func(e *GitHubIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeGitHubOptions(options OptionChanges) func(*GitHubIDPChangedEvent) {
	return func(e *GitHubIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func ChangeGitHubScopes(scopes []string) func(*GitHubIDPChangedEvent) {
	return func(e *GitHubIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func (e *GitHubIDPChangedEvent) Payload() any {
	return e
}

func (e *GitHubIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GitHubIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GitHubIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-Sfrth", "unable to unmarshal event")
	}

	return e, nil
}

type GitHubEnterpriseIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                    string              `json:"id"`
	Name                  string              `json:"name,omitempty"`
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	UserEndpoint          string              `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	Options
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
		*base,
		id,
		name,
		clientID,
		clientSecret,
		authorizationEndpoint,
		tokenEndpoint,
		userEndpoint,
		scopes,
		options,
	}
}

func (e *GitHubEnterpriseIDPAddedEvent) Payload() any {
	return e
}

func (e *GitHubEnterpriseIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GitHubEnterpriseIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GitHubEnterpriseIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-sdfs3", "unable to unmarshal event")
	}

	return e, nil
}

type GitHubEnterpriseIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                    string              `json:"id"`
	Name                  *string             `json:"name,omitempty"`
	ClientID              *string             `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint *string             `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         *string             `json:"tokenEndpoint,omitempty"`
	UserEndpoint          *string             `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewGitHubEnterpriseIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GitHubEnterpriseIDPChanges,
) (*GitHubEnterpriseIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-JHKs9", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubEnterpriseIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type GitHubEnterpriseIDPChanges func(*GitHubEnterpriseIDPChangedEvent)

func ChangeGitHubEnterpriseName(name string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.Name = &name
	}
}
func ChangeGitHubEnterpriseClientID(clientID string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeGitHubEnterpriseClientSecret(clientSecret *crypto.CryptoValue) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeGitHubEnterpriseOptions(options OptionChanges) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func ChangeGitHubEnterpriseAuthorizationEndpoint(authorizationEndpoint string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.AuthorizationEndpoint = &authorizationEndpoint
	}
}

func ChangeGitHubEnterpriseTokenEndpoint(tokenEndpoint string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.TokenEndpoint = &tokenEndpoint
	}
}

func ChangeGitHubEnterpriseUserEndpoint(userEndpoint string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.UserEndpoint = &userEndpoint
	}
}

func ChangeGitHubEnterpriseScopes(scopes []string) func(*GitHubEnterpriseIDPChangedEvent) {
	return func(e *GitHubEnterpriseIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func (e *GitHubEnterpriseIDPChangedEvent) Payload() any {
	return e
}

func (e *GitHubEnterpriseIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GitHubEnterpriseIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GitHubEnterpriseIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-ASf3r", "unable to unmarshal event")
	}

	return e, nil
}
