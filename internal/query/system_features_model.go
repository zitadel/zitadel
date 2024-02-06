package query

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type SystemFeaturesReadModel struct {
	*eventstore.ReadModel
	defaults *feature.Features
	system   *SystemFeatures
}

func NewSystemFeaturesReadModel(defaults *feature.Features) *SystemFeaturesReadModel {
	m := &SystemFeaturesReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   "SYSTEM",
			ResourceOwner: "SYSTEM",
		},
		defaults: defaults,
		system:   new(SystemFeatures),
	}
	m.populateFromDefaults()
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
			feature_v2.SystemDefaultLoginInstanceEventType,
			feature_v2.SystemTriggerIntrospectionProjectionsEventType,
			feature_v2.SystemLegacyIntrospectionEventType,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *SystemFeaturesReadModel) reduceReset() {
	if m.populateFromDefaults() {
		return
	}
	m.system.LoginDefaultOrg = FeatureSource[bool]{}
	m.system.TriggerIntrospectionProjections = FeatureSource[bool]{}
	m.system.LegacyIntrospection = FeatureSource[bool]{}
}

func (m *SystemFeaturesReadModel) populateFromDefaults() bool {
	if m.defaults == nil {
		return false
	}
	m.system.LoginDefaultOrg = FeatureSource[bool]{
		Level: feature.LevelDefault,
		Value: m.defaults.LoginDefaultOrg,
	}
	m.system.TriggerIntrospectionProjections = FeatureSource[bool]{
		Level: feature.LevelDefault,
		Value: m.defaults.TriggerIntrospectionProjections,
	}
	m.system.LegacyIntrospection = FeatureSource[bool]{
		Level: feature.LevelDefault,
		Value: m.defaults.LegacyIntrospection,
	}
	return true
}

func (m *SystemFeaturesReadModel) reduceBoolFeature(event *feature_v2.SetEvent[bool]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	var dst *FeatureSource[bool]

	switch key {
	case feature.KeyLoginDefaultOrg:
		dst = &m.system.LoginDefaultOrg
	case feature.KeyTriggerIntrospectionProjections:
		dst = &m.system.TriggerIntrospectionProjections
	case feature.KeyLegacyIntrospection:
		dst = &m.system.LegacyIntrospection
	}

	*dst = FeatureSource[bool]{
		Level: level,
		Value: event.Value,
	}
	return nil
}
