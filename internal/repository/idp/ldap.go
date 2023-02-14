package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type LDAPIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Host                string              `json:"host"`
	Port                string              `json:"port,omitempty"`
	Tls                 bool                `json:"tls"`
	BaseDN              string              `json:"baseDN"`
	UserObjectClass     string              `json:"userObjectClass"`
	UserUniqueAttribute string              `json:"userUniqueAttribute"`
	Admin               string              `json:"admin"`
	Password            *crypto.CryptoValue `json:"password"`

	IDAttribute                string `json:"idAttribute,omitempty"`
	FirstNameAttribute         string `json:"firstNameAttribute,omitempty"`
	LastNameAttribute          string `json:"lastNameAttribute,omitempty"`
	DisplayNameAttribute       string `json:"displayNameAttribute,omitempty"`
	NickNameAttribute          string `json:"nickNameAttribute,omitempty"`
	PreferredUsernameAttribute string `json:"preferredUsernameAttribute,omitempty"`
	EmailAttribute             string `json:"emailAttribute,omitempty"`
	EmailVerifiedAttribute     string `json:"emailVerifiedAttribute,omitempty"`
	PhoneAttribute             string `json:"phoneAttribute,omitempty"`
	PhoneVerifiedAttribute     string `json:"phoneVerifiedAttribute,omitempty"`
	PreferredLanguageAttribute string `json:"preferredLanguageAttribute,omitempty"`
	AvatarURLAttribute         string `json:"avatarURLAttribute,omitempty"`
	ProfileAttribute           string `json:"profileAttribute,omitempty"`

	Options
}

func NewLDAPIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *LDAPIDPAddedEvent {
	return &LDAPIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *LDAPIDPAddedEvent) Data() interface{} {
	return e
}

func (e *LDAPIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type LDAPIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                  string              `json:"id"`
	Name                *string             `json:"name,omitempty"`
	Host                *string             `json:"host,omitempty"`
	Port                *string             `json:"port,omitempty"`
	Tls                 *bool               `json:"tls,omitempty"`
	BaseDN              *string             `json:"baseDN,omitempty"`
	UserObjectClass     *string             `json:"userObjectClass,omitempty"`
	UserUniqueAttribute *string             `json:"userUniqueAttribute,omitempty"`
	Admin               *string             `json:"admin,omitempty"`
	Password            *crypto.CryptoValue `json:"password,omitempty"`

	IDAttribute                *string `json:"idAttribute,omitempty"`
	FirstNameAttribute         *string `json:"firstNameAttribute,omitempty"`
	LastNameAttribute          *string `json:"lastNameAttribute,omitempty"`
	DisplayNameAttribute       *string `json:"displayNameAttribute,omitempty"`
	NickNameAttribute          *string `json:"nickNameAttribute,omitempty"`
	PreferredUsernameAttribute *string `json:"preferredUsernameAttribute,omitempty"`
	EmailAttribute             *string `json:"emailAttribute,omitempty"`
	EmailVerifiedAttribute     *string `json:"emailVerifiedAttribute,omitempty"`
	PhoneAttribute             *string `json:"phoneAttribute,omitempty"`
	PhoneVerifiedAttribute     *string `json:"phoneVerifiedAttribute,omitempty"`
	PreferredLanguageAttribute *string `json:"preferredLanguageAttribute,omitempty"`
	AvatarURLAttribute         *string `json:"avatarURLAttribute,omitempty"`
	ProfileAttribute           *string `json:"profileAttribute,omitempty"`
	OptionChanges
}

func NewLDAPIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []LDAPIDPChanges,
) (*LDAPIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-S2fa1", "Errors.NoChangesFound")
	}
	changedEvent := &LDAPIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type LDAPIDPChanges func(*LDAPIDPChangedEvent)

func ChangeLDAPClientID(clientID string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeLDAPClientSecret(clientSecret *crypto.CryptoValue) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeLDAPScopes(scopes []string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeLDAPOptions(options OptionChanges) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *LDAPIDPChangedEvent) Data() interface{} {
	return e
}

func (e *LDAPIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Df3f2", "unable to unmarshal event")
	}

	return e, nil
}
