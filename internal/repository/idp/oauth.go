package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type OAuthIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                    string              `json:"id"`
	Name                  string              `json:"name,omitempty"`
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	UserEndpoint          string              `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	IDAttribute           string              `json:"idAttribute,omitempty"`
	Options
}

func NewOAuthIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint,
	idAttribute string,
	scopes []string,
	options Options,
) *OAuthIDPAddedEvent {
	return &OAuthIDPAddedEvent{
		BaseEvent:             *base,
		ID:                    id,
		Name:                  name,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		UserEndpoint:          userEndpoint,
		Scopes:                scopes,
		IDAttribute:           idAttribute,
		Options:               options,
	}
}

func (e *OAuthIDPAddedEvent) Data() interface{} {
	return e
}

func (e *OAuthIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OAuthIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OAuthIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Et1dq", "unable to unmarshal event")
	}

	return e, nil
}

type OAuthIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                    string              `json:"id"`
	Name                  *string             `json:"name,omitempty"`
	ClientID              *string             `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint *string             `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         *string             `json:"tokenEndpoint,omitempty"`
	UserEndpoint          *string             `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	IDAttribute           *string             `json:"idAttribute,omitempty"`
	OptionChanges
}

func NewOAuthIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*OAuthIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &OAuthIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type OAuthIDPChanges func(*OAuthIDPChangedEvent)

func ChangeOAuthName(name string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.Name = &name
	}
}
func ChangeOAuthClientID(clientID string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeOAuthClientSecret(clientSecret *crypto.CryptoValue) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeOAuthOptions(options OptionChanges) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func ChangeOAuthAuthorizationEndpoint(authorizationEndpoint string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.AuthorizationEndpoint = &authorizationEndpoint
	}
}

func ChangeOAuthTokenEndpoint(tokenEndpoint string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.TokenEndpoint = &tokenEndpoint
	}
}

func ChangeOAuthUserEndpoint(userEndpoint string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.UserEndpoint = &userEndpoint
	}
}

func ChangeOAuthScopes(scopes []string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeOAuthIDAttribute(idAttribute string) func(*OAuthIDPChangedEvent) {
	return func(e *OAuthIDPChangedEvent) {
		e.IDAttribute = &idAttribute
	}
}

func (e *OAuthIDPChangedEvent) Data() interface{} {
	return e
}

func (e *OAuthIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OAuthIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OAuthIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-SAf3gw", "unable to unmarshal event")
	}

	return e, nil
}
