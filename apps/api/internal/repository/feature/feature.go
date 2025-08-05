// Package feature implements the v1 feature repository.
// DEPRECATED: use ./feature_v2 instead.
package feature

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

var (
	DefaultLoginInstanceEventType = eventTypePrefix + eventstore.EventType(strings.ToLower("FeatureLoginDefaultOrg")) + setSuffix
)

// DefaultLoginInstanceEventToV2 upgrades the SetEvent to a V2 SetEvent so that
// the v2 reducers can handle the V1 events.
func DefaultLoginInstanceEventToV2(e *SetEvent[Boolean]) *feature_v2.SetEvent[bool] {
	v2e := &feature_v2.SetEvent[bool]{
		BaseEvent: e.BaseEvent,
		Value:     e.Value.Boolean,
	}

	// v1 used a random aggregate ID.
	// v2 uses the instance ID as aggregate ID.
	v2e.BaseEvent.Agg.ID = e.Agg.InstanceID
	v2e.BaseEvent.EventType = feature_v2.InstanceLoginDefaultOrgEventType
	return v2e
}

type SetEvent[T SetEventType] struct {
	*eventstore.BaseEvent `json:"-"`

	Value T
}

func (e *SetEvent[T]) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEvent[T]) Payload() interface{} {
	return e
}

func (e *SetEvent[T]) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SetEventType interface {
	Boolean
	FeatureType() domain.FeatureType
}

type EventType[T SetEventType] struct {
	eventstore.EventType
}

type Boolean struct {
	Boolean bool
}

func (b Boolean) FeatureType() domain.FeatureType {
	return domain.FeatureTypeBoolean
}

func NewSetEvent[T SetEventType](
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	eventType eventstore.EventType,
	setType T,
) *SetEvent[T] {
	return &SetEvent[T]{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, eventType),
		setType,
	}
}
