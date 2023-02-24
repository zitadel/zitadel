package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

const (
	OIDCIDPAddedEventType     eventstore.EventType = "instance.idp.oidc.added"
	OIDCIDPChangedEventType   eventstore.EventType = "instance.idp.oidc.changed"
	JWTIDPAddedEventType      eventstore.EventType = "instance.idp.jwt.added"
	JWTIDPChangedEventType    eventstore.EventType = "instance.idp.jwt.changed"
	GoogleIDPAddedEventType   eventstore.EventType = "instance.idp.google.added"
	GoogleIDPChangedEventType eventstore.EventType = "instance.idp.google.changed"
	LDAPIDPAddedEventType     eventstore.EventType = "instance.idp.ldap.added"
	LDAPIDPChangedEventType   eventstore.EventType = "instance.idp.ldap.changed"
	IDPRemovedEventType       eventstore.EventType = "instance.idp.removed"
)

type OIDCIDPAddedEvent struct {
	idp.OIDCIDPAddedEvent
}

func NewOIDCIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *OIDCIDPAddedEvent {

	return &OIDCIDPAddedEvent{
		OIDCIDPAddedEvent: *idp.NewOIDCIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OIDCIDPAddedEventType,
			),
			id,
			name,
			issuer,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func OIDCIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPAddedEvent{OIDCIDPAddedEvent: *e.(*idp.OIDCIDPAddedEvent)}, nil
}

type OIDCIDPChangedEvent struct {
	idp.OIDCIDPChangedEvent
}

func NewOIDCIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.OIDCIDPChanges,
) (*OIDCIDPChangedEvent, error) {

	changedEvent, err := idp.NewOIDCIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OIDCIDPChangedEvent{OIDCIDPChangedEvent: *changedEvent}, nil
}

func OIDCIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPChangedEvent{OIDCIDPChangedEvent: *e.(*idp.OIDCIDPChangedEvent)}, nil
}

type JWTIDPAddedEvent struct {
	idp.JWTIDPAddedEvent
}

func NewJWTIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) *JWTIDPAddedEvent {

	return &JWTIDPAddedEvent{
		JWTIDPAddedEvent: *idp.NewJWTIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				JWTIDPAddedEventType,
			),
			id,
			name,
			issuer,
			jwtEndpoint,
			keysEndpoint,
			headerName,
			options,
		),
	}
}

func JWTIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.JWTIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &JWTIDPAddedEvent{JWTIDPAddedEvent: *e.(*idp.JWTIDPAddedEvent)}, nil
}

type JWTIDPChangedEvent struct {
	idp.JWTIDPChangedEvent
}

func NewJWTIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.JWTIDPChanges,
) (*JWTIDPChangedEvent, error) {

	changedEvent, err := idp.NewJWTIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			JWTIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &JWTIDPChangedEvent{JWTIDPChangedEvent: *changedEvent}, nil
}

func JWTIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.JWTIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &JWTIDPChangedEvent{JWTIDPChangedEvent: *e.(*idp.JWTIDPChangedEvent)}, nil
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
