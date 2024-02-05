package feature_v2

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
)

var (
	SystemResetEventType                           = resetEventTypeFromFeature(feature.LevelSystem)
	SystemDefaultLoginInstanceEventType            = setEventTypeFromFeature(feature.LevelSystem, feature.LoginDefaultOrg)
	SystemTriggerIntrospectionProjectionsEventType = setEventTypeFromFeature(feature.LevelSystem, feature.TriggerIntrospectionProjections)
	SystemLegacyIntrospectionEventType             = setEventTypeFromFeature(feature.LevelSystem, feature.TriggerIntrospectionProjections)

	InstanceResetEventType                           = resetEventTypeFromFeature(feature.LevelInstance)
	InstanceDefaultLoginInstanceEventType            = setEventTypeFromFeature(feature.LevelInstance, feature.LoginDefaultOrg)
	InstanceTriggerIntrospectionProjectionsEventType = setEventTypeFromFeature(feature.LevelInstance, feature.TriggerIntrospectionProjections)
	InstanceLegacyIntrospectionEventType             = setEventTypeFromFeature(feature.LevelInstance, feature.TriggerIntrospectionProjections)
)

const (
	resetSuffix = "reset"
	setSuffix   = "set"
)

func resetEventTypeFromFeature(level feature.Level) eventstore.EventType {
	return eventstore.EventType(strings.Join([]string{AggregateType, level.String(), resetSuffix}, "."))
}

func setEventTypeFromFeature(level feature.Level, feature feature.Feature) eventstore.EventType {
	return eventstore.EventType(strings.Join([]string{AggregateType, level.String(), feature.String(), setSuffix}, "."))
}

type ResetEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *ResetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *ResetEvent) Payload() interface{} {
	return e
}

func (e *ResetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewResetEvent(
	ctx context.Context,
	aggregate *Aggregate,
	eventType eventstore.EventType,
) *ResetEvent {
	return &ResetEvent{
		eventstore.NewBaseEventForPush(
			ctx, &aggregate.Aggregate, eventType),
	}
}

type SetEvent[T any] struct {
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

func NewSetEvent[T any](
	ctx context.Context,
	aggregate *Aggregate,
	eventType eventstore.EventType,
	value T,
) *SetEvent[T] {
	return &SetEvent[T]{
		eventstore.NewBaseEventForPush(
			ctx, &aggregate.Aggregate, eventType),
		value,
	}
}
