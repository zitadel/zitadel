package deviceauth

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix   eventstore.EventType = "device.authorization."
	AddedEventType                         = eventTypePrefix + "added"
	ApprovedEventType                      = eventTypePrefix + "approved"
	CanceledEventType                      = eventTypePrefix + "canceled"
	DoneEventType                          = eventTypePrefix + "done"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ClientID         string
	DeviceCode       string
	UserCode         string
	Expires          time.Time
	Scopes           []string
	Audience         []string
	State            domain.DeviceAuthState
	NeedRefreshToken bool
}

func (e *AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *AddedEvent) Payload() any {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return NewAddUniqueConstraints(e.DeviceCode, e.UserCode)
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID string,
	deviceCode string,
	userCode string,
	expires time.Time,
	scopes []string,
	audience []string,
	needRefreshToken bool,
) *AddedEvent {
	return &AddedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, AddedEventType,
		),
		clientID, deviceCode, userCode, expires, scopes, audience,
		domain.DeviceAuthStateInitiated, needRefreshToken,
	}
}

type ApprovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID            string
	UserOrgID         string
	UserAuthMethods   []domain.UserAuthMethodType
	AuthTime          time.Time
	PreferredLanguage *language.Tag
	UserAgent         *domain.UserAgent
}

func (e *ApprovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *ApprovedEvent) Payload() any {
	return e
}

func (e *ApprovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewApprovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	userOrgID string,
	userAuthMethods []domain.UserAuthMethodType,
	authTime time.Time,
	preferredLanguage *language.Tag,
	userAgent *domain.UserAgent,
) *ApprovedEvent {
	return &ApprovedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, ApprovedEventType,
		),
		userID,
		userOrgID,
		userAuthMethods,
		authTime,
		preferredLanguage,
		userAgent,
	}
}

type CanceledEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Reason domain.DeviceAuthCanceled
}

func (e *CanceledEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *CanceledEvent) Payload() any {
	return e
}

func (e *CanceledEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCanceledEvent(ctx context.Context, aggregate *eventstore.Aggregate, reason domain.DeviceAuthCanceled) *CanceledEvent {
	return &CanceledEvent{eventstore.NewBaseEventForPush(ctx, aggregate, CanceledEventType), reason}
}

type DoneEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *DoneEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *DoneEvent) Payload() any {
	return e
}

func (e *DoneEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDoneEvent(ctx context.Context, aggregate *eventstore.Aggregate) *DoneEvent {
	return &DoneEvent{eventstore.NewBaseEventForPush(ctx, aggregate, DoneEventType)}
}
