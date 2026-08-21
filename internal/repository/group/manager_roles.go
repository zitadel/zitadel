package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupManagerRolesSetEventType = groupEventTypePrefix + "manager.roles.set"
)

// GroupManagerRolesSetEvent sets the ZITADEL manager roles all members of the group
// receive for the group's organization. An empty role list removes them.
type GroupManagerRolesSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles []string `json:"roles"`
}

func (e *GroupManagerRolesSetEvent) Payload() interface{} {
	return e
}

func (e *GroupManagerRolesSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupManagerRolesSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	roles []string,
) *GroupManagerRolesSetEvent {
	return &GroupManagerRolesSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupManagerRolesSetEventType,
		),
		Roles: roles,
	}
}

func GroupManagerRolesSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupManagerRolesSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-mRq2Lw", "unable to unmarshal group manager roles")
	}

	return e, nil
}
