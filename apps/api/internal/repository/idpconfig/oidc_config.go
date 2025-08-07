package idpconfig

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	OIDCConfigAddedEventType   eventstore.EventType = "oidc.config.added"
	OIDCConfigChangedEventType eventstore.EventType = "oidc.config.changed"
)

type OIDCConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID           string              `json:"idpConfigId"`
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer                string              `json:"issuer,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`

	IDPDisplayNameMapping domain.OIDCMappingField `json:"idpDisplayNameMapping,omitempty"`
	UserNameMapping       domain.OIDCMappingField `json:"usernameMapping,omitempty"`
}

func (e *OIDCConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOIDCConfigAddedEvent(
	base *eventstore.BaseEvent,
	clientID,
	idpConfigID,
	issuer,
	authorizationEndpoint,
	tokenEndpoint string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) *OIDCConfigAddedEvent {

	return &OIDCConfigAddedEvent{
		BaseEvent:             *base,
		IDPConfigID:           idpConfigID,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		Issuer:                issuer,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		Scopes:                scopes,
		IDPDisplayNameMapping: idpDisplayNameMapping,
		UserNameMapping:       userNameMapping,
	}
}

func OIDCConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type OIDCConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`

	ClientID              *string             `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer                *string             `json:"issuer,omitempty"`
	AuthorizationEndpoint *string             `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         *string             `json:"tokenEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`

	IDPDisplayNameMapping *domain.OIDCMappingField `json:"idpDisplayNameMapping,omitempty"`
	UserNameMapping       *domain.OIDCMappingField `json:"usernameMapping,omitempty"`
}

func (e *OIDCConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOIDCConfigChangedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	changes []OIDCConfigChanges,
) (*OIDCConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDPCONFIG-ADzr5", "Errors.NoChangesFound")
	}
	changeEvent := &OIDCConfigChangedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OIDCConfigChanges func(*OIDCConfigChangedEvent)

func ChangeClientID(clientID string) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeClientSecret(secret *crypto.CryptoValue) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ClientSecret = secret
	}
}

func ChangeIssuer(issuer string) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeAuthorizationEndpoint(authorizationEndpoint string) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AuthorizationEndpoint = &authorizationEndpoint
	}
}

func ChangeTokenEndpoint(tokenEndpoint string) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.TokenEndpoint = &tokenEndpoint
	}
}

func ChangeIDPDisplayNameMapping(idpDisplayNameMapping domain.OIDCMappingField) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.IDPDisplayNameMapping = &idpDisplayNameMapping
	}
}

func ChangeUserNameMapping(userNameMapping domain.OIDCMappingField) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.UserNameMapping = &userNameMapping
	}
}

func ChangeScopes(scopes []string) func(*OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.Scopes = scopes
	}
}

func OIDCConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
