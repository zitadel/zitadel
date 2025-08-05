package samlsession

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	samlSessionEventPrefix  = "saml_session."
	AddedType               = samlSessionEventPrefix + "added"
	SAMLResponseAddedType   = samlSessionEventPrefix + "saml_response.added"
	SAMLResponseRevokedType = samlSessionEventPrefix + "saml_response.revoked"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID            string                      `json:"userID"`
	UserResourceOwner string                      `json:"userResourceOwner"`
	SessionID         string                      `json:"sessionID"`
	EntityID          string                      `json:"entityID"`
	Audience          []string                    `json:"audience"`
	AuthMethods       []domain.UserAuthMethodType `json:"authMethods"`
	AuthTime          time.Time                   `json:"authTime"`
	PreferredLanguage *language.Tag               `json:"preferredLanguage,omitempty"`
	UserAgent         *domain.UserAgent           `json:"userAgent,omitempty"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *AddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	userResourceOwner,
	sessionID,
	entityID string,
	audience []string,
	authMethods []domain.UserAuthMethodType,
	authTime time.Time,
	preferredLanguage *language.Tag,
	userAgent *domain.UserAgent,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		UserID:            userID,
		UserResourceOwner: userResourceOwner,
		SessionID:         sessionID,
		EntityID:          entityID,
		Audience:          audience,
		AuthMethods:       authMethods,
		AuthTime:          authTime,
		PreferredLanguage: preferredLanguage,
		UserAgent:         userAgent,
	}
}

type SAMLResponseAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID       string        `json:"id,omitempty"`
	Lifetime time.Duration `json:"lifetime,omitempty"`
}

func (e *SAMLResponseAddedEvent) Payload() interface{} {
	return e
}

func (e *SAMLResponseAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SAMLResponseAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewSAMLResponseAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	lifetime time.Duration,
) *SAMLResponseAddedEvent {
	return &SAMLResponseAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLResponseAddedType,
		),
		ID:       id,
		Lifetime: lifetime,
	}
}

type SAMLResponseRevokedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *SAMLResponseRevokedEvent) Payload() interface{} {
	return e
}

func (e *SAMLResponseRevokedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SAMLResponseRevokedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewSAMLResponseRevokedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *SAMLResponseRevokedEvent {
	return &SAMLResponseRevokedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLResponseRevokedType,
		),
	}
}
