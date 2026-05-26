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

	// TODO(legacy-cleanup): Remove after all pre-migration group memberships have been
	// re-added via AddUsersToGroup (creating new group.users.added events) and the
	// old group.user.added / group.user.removed events are no longer in the eventstore.
	// Legacy event types recorded before the batch-user rename. Never write these.
	LegacyGroupUserAddedEventType   = groupEventTypePrefix + "user.added"
	LegacyGroupUserRemovedEventType = groupEventTypePrefix + "user.removed"

	// TODO(legacy-cleanup): Remove after all pre-migration group memberships have been
	// re-added via AddUsersToGroup and the old group.member.added / group.member.removed
	// events are no longer in the eventstore.
	// Legacy event types from an even older version that used "member" terminology.
	LegacyGroupMemberAddedEventType   = groupEventTypePrefix + "member.added"
	LegacyGroupMemberRemovedEventType = groupEventTypePrefix + "member.removed"
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

// TODO(legacy-cleanup): Remove LegacyGroupUserAddedEvent once all pre-migration group
// memberships have been re-added via AddUsersToGroup.
//
// LegacyGroupUserAddedEvent is the old single-user-per-event form of group membership add.
// Payload: {"userId": "<id>"}
// It is only used for backward-compatible event replay and is never written.
type LegacyGroupUserAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserID               string `json:"userId"`
}

func (e *LegacyGroupUserAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *LegacyGroupUserAddedEvent) Payload() interface{} { return nil }

func (e *LegacyGroupUserAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint { return nil }

// TODO(legacy-cleanup): Remove LegacyGroupUserRemovedEvent once all pre-migration group
// memberships have been re-added via AddUsersToGroup.
//
// LegacyGroupUserRemovedEvent is the old single-user-per-event form of group membership remove.
// Payload: {"userId": "<id>"}
// It is only used for backward-compatible event replay and is never written.
type LegacyGroupUserRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserID               string `json:"userId"`
}

func (e *LegacyGroupUserRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *LegacyGroupUserRemovedEvent) Payload() interface{} { return nil }

func (e *LegacyGroupUserRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint { return nil }

// TODO(legacy-cleanup): Remove LegacyGroupMemberAddedEvent once all pre-migration group
// memberships have been re-added via AddUsersToGroup.
//
// LegacyGroupMemberAddedEvent is the oldest single-user-per-event form using "member" terminology.
// Payload: {"userId": "<id>"}
// It is only used for backward-compatible event replay and is never written.
type LegacyGroupMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserID               string `json:"userId"`
}

func (e *LegacyGroupMemberAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *LegacyGroupMemberAddedEvent) Payload() interface{} { return nil }

func (e *LegacyGroupMemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint { return nil }

// TODO(legacy-cleanup): Remove LegacyGroupMemberRemovedEvent once all pre-migration group
// memberships have been re-added via AddUsersToGroup.
//
// LegacyGroupMemberRemovedEvent is the oldest single-user-per-event form using "member" terminology.
// Payload: {"userId": "<id>"}
// It is only used for backward-compatible event replay and is never written.
type LegacyGroupMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserID               string `json:"userId"`
}

func (e *LegacyGroupMemberRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *LegacyGroupMemberRemovedEvent) Payload() interface{} { return nil }

func (e *LegacyGroupMemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// UnmarshalJSON supports two on-disk formats for backward compatibility:
//
//   - New format (current): {"users":[{"userId":"…","attributes":{…}},...]}
//   - Old format (pre-attributes): {"userIds":["id1","id2",...]}
//
// The old format was used before per-user attributes were introduced.
// Old events are re-hydrated into the new struct with empty attribute maps.
func (e *GroupUsersAddedEvent) UnmarshalJSON(data []byte) error {
	// Try new format first
	type newFormat struct {
		Users []GroupUser `json:"users"`
	}
	var nf newFormat
	if err := json.Unmarshal(data, &nf); err == nil && len(nf.Users) > 0 {
		e.Users = nf.Users
		return nil
	}

	// Fall back to old format: {"userIds":["id1","id2"]}
	type oldFormat struct {
		UserIDs []string `json:"userIds"`
	}
	var of oldFormat
	if err := json.Unmarshal(data, &of); err != nil {
		return err
	}
	e.Users = make([]GroupUser, len(of.UserIDs))
	for i, id := range of.UserIDs {
		e.Users[i] = GroupUser{UserID: id}
	}
	return nil
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
