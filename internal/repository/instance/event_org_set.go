package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DefaultOrgSetEventType eventstore.EventType = "instance.default.org.set"
)

type DefaultOrgSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	OrgID string `json:"orgId"`
}

func (e *DefaultOrgSetEvent) Payload() interface{} {
	return e
}

func (e *DefaultOrgSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDefaultOrgSetEventEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	orgID string,
) *DefaultOrgSetEvent {
	return &DefaultOrgSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DefaultOrgSetEventType,
		),
		OrgID: orgID,
	}
}

func DefaultOrgSetMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DefaultOrgSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-cdFZH", "unable to unmarshal default org set")
	}

	return e, nil
}
