package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type OAuthIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                    string              `json:"id"`
	Name                  string              `json:"name,omitempty"`
	ClientID              string              `json:"client_id,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"client_secret,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	UserEndpoint          string              `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
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
	userEndpoint string,
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
		Options:               options,
	}
}

func (e *OAuthIDPAddedEvent) Data() interface{} {
	return e
}

func (e *OAuthIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
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

	oldName string

	ID                    string              `json:"id"`
	Name                  *string             `json:"name,omitempty"`
	ClientID              *string             `json:"client_id,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"client_secret,omitempty"`
	AuthorizationEndpoint *string             `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         *string             `json:"tokenEndpoint,omitempty"`
	UserEndpoint          *string             `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"` // TODO: tristate?
	OptionChanges
}

func NewOAuthIDPChangedEvent(
	base *eventstore.BaseEvent,
	id,
	oldName string,
	changes []OAuthIDPChanges,
) (*OAuthIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &OAuthIDPChangedEvent{
		BaseEvent: *base,
		oldName:   oldName,
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

func (e *OAuthIDPChangedEvent) Data() interface{} {
	return e
}

func (e *OAuthIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	if e.Name == nil || e.oldName == *e.Name { // TODO: nil check should be enough
		return nil
	}
	return []*eventstore.EventUniqueConstraint{
		idpconfig.NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
		idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	}
}

func OAuthIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OAuthIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
	}

	return e, nil
}
