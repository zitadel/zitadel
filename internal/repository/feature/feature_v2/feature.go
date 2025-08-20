package feature_v2

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	SystemResetEventType                          = resetEventTypeFromFeature(feature.LevelSystem)
	SystemLoginDefaultOrgEventType                = setEventTypeFromFeature(feature.LevelSystem, feature.KeyLoginDefaultOrg)
	SystemUserSchemaEventType                     = setEventTypeFromFeature(feature.LevelSystem, feature.KeyUserSchema)
	SystemTokenExchangeEventType                  = setEventTypeFromFeature(feature.LevelSystem, feature.KeyTokenExchange)
	SystemImprovedPerformanceEventType            = setEventTypeFromFeature(feature.LevelSystem, feature.KeyImprovedPerformance)
	SystemOIDCSingleV1SessionTerminationEventType = setEventTypeFromFeature(feature.LevelSystem, feature.KeyOIDCSingleV1SessionTermination)
	SystemEnableBackChannelLogout                 = setEventTypeFromFeature(feature.LevelSystem, feature.KeyEnableBackChannelLogout)
	SystemLoginVersion                            = setEventTypeFromFeature(feature.LevelSystem, feature.KeyLoginV2)
	SystemPermissionCheckV2                       = setEventTypeFromFeature(feature.LevelSystem, feature.KeyPermissionCheckV2)

	InstanceResetEventType                          = resetEventTypeFromFeature(feature.LevelInstance)
	InstanceLoginDefaultOrgEventType                = setEventTypeFromFeature(feature.LevelInstance, feature.KeyLoginDefaultOrg)
	InstanceUserSchemaEventType                     = setEventTypeFromFeature(feature.LevelInstance, feature.KeyUserSchema)
	InstanceTokenExchangeEventType                  = setEventTypeFromFeature(feature.LevelInstance, feature.KeyTokenExchange)
	InstanceImprovedPerformanceEventType            = setEventTypeFromFeature(feature.LevelInstance, feature.KeyImprovedPerformance)
	InstanceDebugOIDCParentErrorEventType           = setEventTypeFromFeature(feature.LevelInstance, feature.KeyDebugOIDCParentError)
	InstanceOIDCSingleV1SessionTerminationEventType = setEventTypeFromFeature(feature.LevelInstance, feature.KeyOIDCSingleV1SessionTermination)
	InstanceEnableBackChannelLogout                 = setEventTypeFromFeature(feature.LevelInstance, feature.KeyEnableBackChannelLogout)
	InstanceLoginVersion                            = setEventTypeFromFeature(feature.LevelInstance, feature.KeyLoginV2)
	InstancePermissionCheckV2                       = setEventTypeFromFeature(feature.LevelInstance, feature.KeyPermissionCheckV2)
	InstanceConsoleUseV2UserApi                     = setEventTypeFromFeature(feature.LevelInstance, feature.KeyConsoleUseV2UserApi)
)

const (
	resetSuffix = "reset"
	setSuffix   = "set"
)

func resetEventTypeFromFeature(level feature.Level) eventstore.EventType {
	return eventstore.EventType(strings.Join([]string{AggregateType, level.String(), resetSuffix}, "."))
}

func setEventTypeFromFeature(level feature.Level, key feature.Key) eventstore.EventType {
	return eventstore.EventType(strings.Join([]string{AggregateType, level.String(), key.String(), setSuffix}, "."))
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

type FeatureJSON struct {
	Key   feature.Key
	Value []byte
}

// FeatureJSON prepares converts the event to a key-value pair with a JSON value payload.
func (e *SetEvent[T]) FeatureJSON() (*FeatureJSON, error) {
	_, key, err := e.FeatureInfo()
	if err != nil {
		return nil, err
	}
	jsonValue, err := json.Marshal(e.Value)
	if err != nil {
		return nil, zerrors.ThrowInternalf(err, "FEAT-go9Ji", "reduce.wrong.event.type %s", e.EventType)
	}
	return &FeatureJSON{
		Key:   key,
		Value: jsonValue,
	}, nil
}

// FeatureInfo extracts a feature's level and key from the event.
func (e *SetEvent[T]) FeatureInfo() (feature.Level, feature.Key, error) {
	ss := strings.Split(string(e.EventType), ".")
	if len(ss) != 4 {
		return 0, 0, zerrors.ThrowInternalf(nil, "FEAT-Ahs4m", "reduce.wrong.event.type %s", e.EventType)
	}
	level, err := feature.LevelString(ss[1])
	if err != nil {
		return 0, 0, zerrors.ThrowInternalf(err, "FEAT-Boo2i", "reduce.wrong.event.type %s", e.EventType)
	}
	key, err := feature.KeyString(ss[2])
	if err != nil {
		return 0, 0, zerrors.ThrowInternalf(err, "FEAT-eir0M", "reduce.wrong.event.type %s", e.EventType)
	}
	return level, key, nil
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
