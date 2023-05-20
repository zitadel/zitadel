package instance

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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

func (e *DefaultOrgSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func DefaultOrgSetMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DefaultOrgSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-cdFZH", "unable to unmarshal default org set")
	}

	return e, nil
}
