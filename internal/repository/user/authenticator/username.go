package authenticator

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	usernamePrefix      = eventPrefix + "username."
	UsernameCreatedType = usernamePrefix + "created"
	UsernameDeletedType = usernamePrefix + "deleted"

	uniqueUsernameType = "username"
)

func NewAddUsernameUniqueConstraint(resourceOwner string, isOrgSpecific bool, username string) *eventstore.UniqueConstraint {
	if isOrgSpecific {
		return eventstore.NewAddEventUniqueConstraint(
			uniqueUsernameType,
			resourceOwner+":"+username,
			"TODO")
	}

	return eventstore.NewAddEventUniqueConstraint(
		uniqueUsernameType,
		username,
		"TODO")
}

func NewRemoveUsernameUniqueConstraint(resourceOwner string, isOrgSpecific bool, username string) *eventstore.UniqueConstraint {
	if isOrgSpecific {
		return eventstore.NewRemoveUniqueConstraint(
			uniqueUsernameType,
			resourceOwner+":"+username,
		)
	}

	return eventstore.NewRemoveUniqueConstraint(
		uniqueUsernameType,
		username,
	)
}

type UsernameCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID        string `json:"userID"`
	IsOrgSpecific bool   `json:"isOrgSpecific"`
	Username      string `json:"username"`
}

func (e *UsernameCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UsernameCreatedEvent) Payload() interface{} {
	return e
}

func (e *UsernameCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUsernameUniqueConstraint(e.Agg.ResourceOwner, e.IsOrgSpecific, e.Username)}
}

func NewUsernameCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	isOrgSpecific bool,
	username string,
) *UsernameCreatedEvent {
	return &UsernameCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UsernameCreatedType,
		),
		UserID:        userID,
		IsOrgSpecific: isOrgSpecific,
		Username:      username,
	}
}

type UsernameDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	isOrgSpecific bool
	username      string
}

func (e *UsernameDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *UsernameDeletedEvent) Payload() interface{} {
	return e
}

func (e *UsernameDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveUsernameUniqueConstraint(e.Agg.ResourceOwner, e.isOrgSpecific, e.username),
	}
}

func NewUsernameDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	isOrgSpecific bool,
	username string,
) *UsernameDeletedEvent {
	return &UsernameDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UsernameDeletedType,
		),
		isOrgSpecific: isOrgSpecific,
		username:      username,
	}
}
