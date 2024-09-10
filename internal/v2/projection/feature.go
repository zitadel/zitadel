package projection

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/feature"
)

type Feature[T any] struct {
	projection

	Level feature.Level
	Key   feature.Key
	value T
}

func NewFeature[T any](level feature.Level, key feature.Key) *Feature[T] {
	return &Feature[T]{
		Level: level,
		Key:   key,
	}
}

func (f *Feature[T]) Value() T {
	return f.value
}

func (f *Feature[T]) Reducers() map[string]map[string]eventstore.ReduceEvent {
	if f.reducers != nil {
		return f.reducers
	}

	f.reducers = map[string]map[string]eventstore.ReduceEvent{
		feature.AggregateType: {
			feature.ResetEventType(f.Level):      f.reduceFeatureReset,
			feature.SetEventType(f.Level, f.Key): f.reduceFeatureSet,
		},
	}

	return f.reducers
}

func (f *Feature[T]) reduceFeatureReset(event *eventstore.StorageEvent) error {
	if !f.ShouldReduce(event) {
		return nil
	}

	e, err := feature.ResetEventFromStorage(event)
	if err != nil {
		return err
	}

	if e.Level != f.Level {
		return nil
	}

	var t T
	f.value = t

	return nil
}

func (f *Feature[T]) reduceFeatureSet(event *eventstore.StorageEvent) error {
	if !f.ShouldReduce(event) {
		return nil
	}

	e, err := feature.SetEventFromStorage(event)
	if err != nil {
		return err
	}

	if e.Level != f.Level || e.Key != f.Key {
		return nil
	}

	return json.Unmarshal(e.Payload.Value, &f.value)
}
