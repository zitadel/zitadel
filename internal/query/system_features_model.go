package query

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type SystemFeaturesReadModel struct {
	*eventstore.ReadModel
	system *SystemFeatures
}

func NewSystemFeaturesReadModel() *SystemFeaturesReadModel {
	m := &SystemFeaturesReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   "SYSTEM",
			ResourceOwner: "SYSTEM",
		},
		system: new(SystemFeatures),
	}
	return m
}

func (m *SystemFeaturesReadModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *feature_v2.ResetEvent:
			m.reduceReset()
		case *feature_v2.SetEvent[bool]:
			err := m.reduceBoolFeature(e)
			if err != nil {
				return err
			}
		case *feature_v2.SetEvent[[]int32]:
			err := m.reduceEnumListFeature(e)
			if err != nil {
				return err
			}
		}
	}
	return m.ReadModel.Reduce()
}

func (m *SystemFeaturesReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AddQuery().
		AggregateTypes(feature_v2.AggregateType).
		AggregateIDs(m.AggregateID).
		EventTypes(
			feature_v2.SystemResetEventType,
			feature_v2.SystemLoginDefaultOrgEventType,
			feature_v2.SystemTriggerIntrospectionProjectionsEventType,
			feature_v2.SystemLegacyIntrospectionEventType,
			feature_v2.SystemUserSchemaEventType,
			feature_v2.SystemTokenExchangeEventType,
			feature_v2.SystemActionsEventType,
			feature_v2.SystemImprovedPerformanceEventType,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *SystemFeaturesReadModel) reduceReset() {
	m.system = new(SystemFeatures)
}

func reduceFeatureV2[T any](features *SystemFeatures, event *feature_v2.SetEvent[T]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyLoginDefaultOrg:
		setSource(level, &features.LoginDefaultOrg, any(event.Value).(bool))
	case feature.KeyTriggerIntrospectionProjections:
		dst = &features.TriggerIntrospectionProjections
	case feature.KeyLegacyIntrospection:
		dst = &features.LegacyIntrospection
	case feature.KeyUserSchema:
		dst = &features.UserSchema
	case feature.KeyTokenExchange:
		dst = &features.TokenExchange
	case feature.KeyActions:
		dst = &features.Actions
	case feature.KeyImprovedPerformance:
		dst = &features.ImprovedPerformance
	}
}

func setSource[T any](level feature.Level, source *FeatureSource[T], value T) {
	source.Level = level
	source.Value = value
}

func (m *SystemFeaturesReadModel) reduceBoolFeature(event *feature_v2.SetEvent[bool]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	var dst *FeatureSource[bool]

	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyLoginDefaultOrg:
		dst = &m.system.LoginDefaultOrg
	case feature.KeyTriggerIntrospectionProjections:
		dst = &m.system.TriggerIntrospectionProjections
	case feature.KeyLegacyIntrospection:
		dst = &m.system.LegacyIntrospection
	case feature.KeyUserSchema:
		dst = &m.system.UserSchema
	case feature.KeyTokenExchange:
		dst = &m.system.TokenExchange
	case feature.KeyActions:
		dst = &m.system.Actions
	}

	*dst = FeatureSource[bool]{
		Level: level,
		Value: event.Value,
	}
	return nil
}

func (m *SystemFeaturesReadModel) reduceEnumListFeature(event *feature_v2.SetEvent[[]int32]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	var dst *FeatureSource[[]int32]

	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyImprovedPerformance:
		dst = &m.system.ImprovedPerformance
	}

	*dst = FeatureSource[[]int32]{
		Level: level,
		Value: event.Value,
	}
	return nil
}
