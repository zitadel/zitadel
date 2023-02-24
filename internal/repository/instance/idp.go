package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

const (
	OAuthIDPAddedEventType    eventstore.EventType = "instance.idp.oauth.added"
	OAuthIDPChangedEventType  eventstore.EventType = "instance.idp.oauth.changed"
	GoogleIDPAddedEventType   eventstore.EventType = "instance.idp.google.added"
	GoogleIDPChangedEventType eventstore.EventType = "instance.idp.google.changed"
	LDAPIDPAddedEventType     eventstore.EventType = "instance.idp.ldap.added"
	LDAPIDPChangedEventType   eventstore.EventType = "instance.idp.ldap.changed"
	IDPRemovedEventType       eventstore.EventType = "instance.idp.removed"
)

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
	id string,
	changes []idp.OAuthIDPChanges,
) (*OAuthIDPChangedEvent, error) {

	changedEvent, err := idp.NewOAuthIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OAuthIDPChangedEventType,
		),
		id,
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

type GoogleIDPAddedEvent struct {
	idp.GoogleIDPAddedEvent
}

func NewGoogleIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
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
			name,
			clientID,
			clientSecret,
			scopes,
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

type LDAPIDPAddedEvent struct {
	idp.LDAPIDPAddedEvent
}

func NewLDAPIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password *crypto.CryptoValue,
	attributes idp.LDAPAttributes,
	options idp.Options,
) *LDAPIDPAddedEvent {

	return &LDAPIDPAddedEvent{
		LDAPIDPAddedEvent: *idp.NewLDAPIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LDAPIDPAddedEventType,
			),
			id,
			name,
			host,
			port,
			tls,
			baseDN,
			userObjectClass,
			userUniqueAttribute,
			admin,
			password,
			attributes,
			options,
		),
	}
}

func LDAPIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.LDAPIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LDAPIDPAddedEvent{LDAPIDPAddedEvent: *e.(*idp.LDAPIDPAddedEvent)}, nil
}

type LDAPIDPChangedEvent struct {
	idp.LDAPIDPChangedEvent
}

func NewLDAPIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName string,
	changes []idp.LDAPIDPChanges,
) (*LDAPIDPChangedEvent, error) {

	changedEvent, err := idp.NewLDAPIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LDAPIDPChangedEventType,
		),
		id,
		oldName,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LDAPIDPChangedEvent{LDAPIDPChangedEvent: *changedEvent}, nil
}

func LDAPIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.LDAPIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LDAPIDPChangedEvent{LDAPIDPChangedEvent: *e.(*idp.LDAPIDPChangedEvent)}, nil
}

type IDPRemovedEvent struct {
	idp.RemovedEvent
}

func NewIDPRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	name string,
) *IDPRemovedEvent {
	return &IDPRemovedEvent{
		RemovedEvent: *idp.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPRemovedEventType,
			),
			id,
			name,
		),
	}
}

func (e *IDPRemovedEvent) Data() interface{} {
	return e
}

func IDPRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPRemovedEvent{RemovedEvent: *e.(*idp.RemovedEvent)}, nil
}
