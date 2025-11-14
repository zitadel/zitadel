package sessionlogout

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix                 = "session_logout."
	backChannelEventTypePrefix      = eventTypePrefix + "back_channel."
	BackChannelLogoutRegisteredType = backChannelEventTypePrefix + "registered"
	BackChannelLogoutSentType       = backChannelEventTypePrefix + "sent"
)

type BackChannelLogoutRegisteredEvent struct {
	*eventstore.BaseEvent `json:"-"`

	OIDCSessionID        string `json:"oidc_session_id"`
	UserID               string `json:"user_id"`
	ClientID             string `json:"client_id"`
	BackChannelLogoutURI string `json:"back_channel_logout_uri"`
}

// Payload implements eventstore.Command.
func (e *BackChannelLogoutRegisteredEvent) Payload() any {
	return e
}

func (e *BackChannelLogoutRegisteredEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *BackChannelLogoutRegisteredEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewBackChannelLogoutRegisteredEvent(ctx context.Context, aggregate *eventstore.Aggregate, oidcSessionID, userID, clientID, backChannelLogoutURI string) *BackChannelLogoutRegisteredEvent {
	return &BackChannelLogoutRegisteredEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			BackChannelLogoutRegisteredType,
		),
		OIDCSessionID:        oidcSessionID,
		UserID:               userID,
		ClientID:             clientID,
		BackChannelLogoutURI: backChannelLogoutURI,
	}
}

type BackChannelLogoutSentEvent struct {
	eventstore.BaseEvent `json:"-"`

	OIDCSessionID string `json:"oidc_session_id"`
}

func (e *BackChannelLogoutSentEvent) Payload() interface{} {
	return e
}

func (e *BackChannelLogoutSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *BackChannelLogoutSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewBackChannelLogoutSentEvent(ctx context.Context, aggregate *eventstore.Aggregate, oidcSessionID string) *BackChannelLogoutSentEvent {
	return &BackChannelLogoutSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			BackChannelLogoutSentType,
		),
		OIDCSessionID: oidcSessionID,
	}
}
