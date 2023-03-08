package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

const (
	OAuthIDPAddedEventType              eventstore.EventType = "instance.idp.oauth.added"
	OAuthIDPChangedEventType            eventstore.EventType = "instance.idp.oauth.changed"
	OIDCIDPAddedEventType               eventstore.EventType = "instance.idp.oidc.added"
	OIDCIDPChangedEventType             eventstore.EventType = "instance.idp.oidc.changed"
	JWTIDPAddedEventType                eventstore.EventType = "instance.idp.jwt.added"
	JWTIDPChangedEventType              eventstore.EventType = "instance.idp.jwt.changed"
	GitHubIDPAddedEventType             eventstore.EventType = "instance.idp.github.added"
	GitHubIDPChangedEventType           eventstore.EventType = "instance.idp.github.changed"
	GitHubEnterpriseIDPAddedEventType   eventstore.EventType = "instance.idp.github_enterprise.added"
	GitHubEnterpriseIDPChangedEventType eventstore.EventType = "instance.idp.github_enterprise.changed"
	GoogleIDPAddedEventType             eventstore.EventType = "instance.idp.google.added"
	GoogleIDPChangedEventType           eventstore.EventType = "instance.idp.google.changed"
	LDAPIDPAddedEventType               eventstore.EventType = "instance.idp.ldap.added"
	LDAPIDPChangedEventType             eventstore.EventType = "instance.idp.ldap.changed"
	IDPRemovedEventType                 eventstore.EventType = "instance.idp.removed"
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
	userEndpoint,
	idAttribute string,
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
			idAttribute,
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

type GitHubIDPAddedEvent struct {
	idp.GitHubIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
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
			name,
			clientID,
			clientSecret,
			scopes,
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
	changes []idp.GitHubIDPChanges,
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

type GitHubEnterpriseIDPAddedEvent struct {
	idp.GitHubEnterpriseIDPAddedEvent
}

func NewGitHubEnterpriseIDPAddedEvent(
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
) *GitHubEnterpriseIDPAddedEvent {

	return &GitHubEnterpriseIDPAddedEvent{
		GitHubEnterpriseIDPAddedEvent: *idp.NewGitHubEnterpriseIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitHubEnterpriseIDPAddedEventType,
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

func GitHubEnterpriseIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GitHubEnterpriseIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPAddedEvent{GitHubEnterpriseIDPAddedEvent: *e.(*idp.GitHubEnterpriseIDPAddedEvent)}, nil
}

type GitHubEnterpriseIDPChangedEvent struct {
	idp.GitHubEnterpriseIDPChangedEvent
}

func NewGitHubEnterpriseIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GitHubEnterpriseIDPChanges,
) (*GitHubEnterpriseIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitHubEnterpriseIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitHubEnterpriseIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitHubEnterpriseIDPChangedEvent{GitHubEnterpriseIDPChangedEvent: *changedEvent}, nil
}

func GitHubEnterpriseIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := idp.GitHubEnterpriseIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPChangedEvent{GitHubEnterpriseIDPChangedEvent: *e.(*idp.GitHubEnterpriseIDPChangedEvent)}, nil
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
) *IDPRemovedEvent {
	return &IDPRemovedEvent{
		RemovedEvent: *idp.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPRemovedEventType,
			),
			id,
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
