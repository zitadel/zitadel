package group

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupUsersAddedEventType   = groupEventTypePrefix + "users.added"
	GroupUsersChangedEventType = groupEventTypePrefix + "users.changed"
	GroupUsersRemovedEventType = groupEventTypePrefix + "users.removed"
)

// AttributeValue supports both single string and string array values in JSON.
// Single strings are normalized to single-element arrays.
type AttributeValue []string

func (av *AttributeValue) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*av = AttributeValue(arr)
		return nil
	}

	// Try single string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return zerrors.ThrowInvalidArgument(err, "GROUP-aK9s3", "attribute value must be string or array")
	}
	*av = AttributeValue{str}
	return nil
}

func (av AttributeValue) MarshalJSON() ([]byte, error) {
	if len(av) == 1 {
		return json.Marshal(av[0])
	}
	return json.Marshal([]string(av))
}

// GroupUser represents a user's membership in a group along with per-user attributes.
type GroupUser struct {
	UserID     string                    `json:"userId"`
	Attributes map[string]AttributeValue `json:"attributes,omitempty"`
}

type GroupUsersAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Users []GroupUser `json:"users"`
}

func NewGroupUsersAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, users []GroupUser) *GroupUsersAddedEvent {
	return &GroupUsersAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUsersAddedEventType,
		),
		Users: users,
	}
}

func (e *GroupUsersAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUsersAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupUsersAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type GroupUsersRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserIDs []string `json:"userIds"`
}

func NewGroupUsersRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userIDs []string,
) *GroupUsersRemovedEvent {
	return &GroupUsersRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUsersRemovedEventType,
		),
		UserIDs: userIDs,
	}
}

func (e *GroupUsersRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUsersRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupUsersRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type GroupUserChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID     string                    `json:"userId,omitempty"`
	Attributes map[string]AttributeValue `json:"attributes,omitempty"`
}

func (e *GroupUserChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUserChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupUserChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	attributes map[string]AttributeValue,
) *GroupUserChangedEvent {
	return &GroupUserChangedEvent{
		BaseEvent:  *eventstore.NewBaseEventForPush(ctx, aggregate, GroupUsersChangedEventType),
		UserID:     userID,
		Attributes: attributes,
	}
}
