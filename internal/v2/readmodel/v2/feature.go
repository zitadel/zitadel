package readmodel

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	v2_feature "github.com/zitadel/zitadel/internal/v2/feature"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/v2/system"
)

var _ objectManager = (*SystemFeatures)(nil)

type SystemFeatures struct {
	readModel
	object *objectReadModel

	LoginDefaultOrg                 *projection.Feature[bool]
	TriggerIntrospectionProjections *projection.Feature[bool]
	LegacyIntrospection             *projection.Feature[bool]
	UserSchema                      *projection.Feature[bool]
	TokenExchange                   *projection.Feature[bool]
	Actions                         *projection.Feature[bool]
	ImprovedPerformance             *projection.Feature[[]ImprovedPerformanceType]
	OIDCSingleV1SessionTermination  *projection.Feature[bool]
}

func NewSystemFeatures(ctx context.Context, eventStore *eventstore.Eventstore) *SystemFeatures {
	features := &SystemFeatures{
		LoginDefaultOrg:                 projection.NewFeature[bool](feature.LevelSystem, feature.KeyLoginDefaultOrg),
		TriggerIntrospectionProjections: projection.NewFeature[bool](feature.LevelSystem, feature.KeyTriggerIntrospectionProjections),
		LegacyIntrospection:             projection.NewFeature[bool](feature.LevelSystem, feature.KeyLegacyIntrospection),
		UserSchema:                      projection.NewFeature[bool](feature.LevelSystem, feature.KeyUserSchema),
		TokenExchange:                   projection.NewFeature[bool](feature.LevelSystem, feature.KeyTokenExchange),
		Actions:                         projection.NewFeature[bool](feature.LevelSystem, feature.KeyActions),
		ImprovedPerformance:             projection.NewFeature[[]ImprovedPerformanceType](feature.LevelSystem, feature.KeyImprovedPerformance),
		OIDCSingleV1SessionTermination:  projection.NewFeature[bool](feature.LevelSystem, feature.KeyOIDCSingleV1SessionTermination),
	}
	features.object = newObjectReadModel(ctx, features, eventStore)
	features.object.init(ctx)

	return features
}

// EventstoreV3Query implements objectManager.
func (s *SystemFeatures) EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder {
	eventTypes := make([]eventstore.EventType, 0, 11)

	for _, featureEventTypes := range s.Reducers() {
		for eventType := range featureEventTypes {
			eventTypes = append(eventTypes, eventstore.EventType(eventType))
		}
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(system.AggregateInstance).
		AddQuery().
		AggregateTypes(eventstore.AggregateType(v2_feature.AggregateType)).
		EventTypes(eventTypes...).
		Builder()
}

// Name implements objectManager.
func (s *SystemFeatures) Name() string {
	return "system_features"
}

// Reducers implements objectManager.
func (s *SystemFeatures) Reducers() projection.Reducers {
	if s.reducers != nil {
		return s.reducers
	}
	s.reducers = projection.Reducers{
		system.AggregateType: make(map[string]v2_es.ReduceEvent, 8),
	}

	s.reducers = projection.MergeReducers(
		s.LoginDefaultOrg.Reducers(),
		s.TriggerIntrospectionProjections.Reducers(),
		s.LegacyIntrospection.Reducers(),
		s.UserSchema.Reducers(),
		s.TokenExchange.Reducers(),
		s.Actions.Reducers(),
		s.ImprovedPerformance.Reducers(),
		s.OIDCSingleV1SessionTermination.Reducers(),
	)

	return s.reducers
}

type ImprovedPerformanceType = feature.ImprovedPerformanceType

type InstanceFeatures struct {
	readModel
	object *objectReadModel

	instanceID string

	WebKey               *projection.Feature[bool]
	DebugOIDCParentError *projection.Feature[bool]
	InMemoryProjections  *projection.Feature[bool]

	LoginDefaultOrg                 *InheritedFeature[bool]
	TriggerIntrospectionProjections *InheritedFeature[bool]
	LegacyIntrospection             *InheritedFeature[bool]
	UserSchema                      *InheritedFeature[bool]
	TokenExchange                   *InheritedFeature[bool]
	Actions                         *InheritedFeature[bool]
	ImprovedPerformance             *InheritedFeature[[]ImprovedPerformanceType]
	OIDCSingleV1SessionTermination  *InheritedFeature[bool]
}

func NewInstanceFeatures(ctx context.Context, eventStore *eventstore.Eventstore, systemFeatures *SystemFeatures, instanceID string) *InstanceFeatures {
	features := &InstanceFeatures{
		instanceID: instanceID,

		WebKey:               projection.NewFeature[bool](feature.LevelInstance, feature.KeyWebKey),
		DebugOIDCParentError: projection.NewFeature[bool](feature.LevelInstance, feature.KeyDebugOIDCParentError),
		InMemoryProjections:  projection.NewFeature[bool](feature.LevelInstance, feature.KeyInMemoryProjections),

		LoginDefaultOrg:                 newInheritedFeature(systemFeatures.LoginDefaultOrg, feature.LevelInstance, feature.KeyLoginDefaultOrg),
		TriggerIntrospectionProjections: newInheritedFeature(systemFeatures.TriggerIntrospectionProjections, feature.LevelInstance, feature.KeyTriggerIntrospectionProjections),
		LegacyIntrospection:             newInheritedFeature(systemFeatures.LegacyIntrospection, feature.LevelInstance, feature.KeyLegacyIntrospection),
		UserSchema:                      newInheritedFeature(systemFeatures.UserSchema, feature.LevelInstance, feature.KeyUserSchema),
		TokenExchange:                   newInheritedFeature(systemFeatures.TokenExchange, feature.LevelInstance, feature.KeyTokenExchange),
		Actions:                         newInheritedFeature(systemFeatures.Actions, feature.LevelInstance, feature.KeyActions),
		ImprovedPerformance:             newInheritedFeature(systemFeatures.ImprovedPerformance, feature.LevelInstance, feature.KeyImprovedPerformance),
		OIDCSingleV1SessionTermination:  newInheritedFeature(systemFeatures.OIDCSingleV1SessionTermination, feature.LevelInstance, feature.KeyOIDCSingleV1SessionTermination),
	}
	features.object = newObjectReadModel(ctx, features, eventStore)
	features.object.init(ctx)

	return features
}

// EventstoreV3Query implements objectManager.
func (i *InstanceFeatures) EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder {
	eventTypes := make([]eventstore.EventType, 0, 11)

	for _, featureEventTypes := range i.Reducers() {
		for eventType := range featureEventTypes {
			eventTypes = append(eventTypes, eventstore.EventType(eventType))
		}
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(i.instanceID).
		AddQuery().
		AggregateTypes(eventstore.AggregateType(v2_feature.AggregateType)).
		EventTypes(eventTypes...).
		Builder()
}

// Name implements objectManager.
func (i *InstanceFeatures) Name() string {
	return "instance_features"
}

// Reducers implements objectManager.
func (f *InstanceFeatures) Reducers() projection.Reducers {
	if f.reducers != nil {
		return f.reducers
	}

	f.reducers = projection.MergeReducers(
		f.WebKey.Reducers(),
		f.DebugOIDCParentError.Reducers(),
		f.InMemoryProjections.Reducers(),
		f.LoginDefaultOrg.Reducers(),
		f.TriggerIntrospectionProjections.Reducers(),
		f.LegacyIntrospection.Reducers(),
		f.UserSchema.Reducers(),
		f.TokenExchange.Reducers(),
		f.Actions.Reducers(),
		f.ImprovedPerformance.Reducers(),
		f.OIDCSingleV1SessionTermination.Reducers(),
	)

	return f.reducers
}

type featureValuer[T any] interface {
	Value() T
}

func newInheritedFeature[T any](parent featureValuer[T], level feature.Level, key feature.Key) *InheritedFeature[T] {
	return &InheritedFeature[T]{
		Feature: projection.NewFeature[*T](level, key),
		parent:  parent,
	}
}

type InheritedFeature[T any] struct {
	*projection.Feature[*T]
	parent featureValuer[T]
}

func (f *InheritedFeature[T]) Value() T {
	if f.Feature.Value() != nil {
		return *f.Feature.Value()
	}

	return f.parent.Value()
}
