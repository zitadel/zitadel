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
			err := reduceSystemFeatureSet(m.system, e)
			if err != nil {
				return err
			}
		case *feature_v2.SetEvent[[]feature.ImprovedPerformanceType]:
			err := reduceSystemFeatureSet(m.system, e)
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
			feature_v2.SystemOIDCSingleV1SessionTerminationEventType,
			feature_v2.SystemDisableUserTokenEvent,
			feature_v2.SystemEnableBackChannelLogout,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *SystemFeaturesReadModel) reduceReset() {
	m.system = nil
	m.system = new(SystemFeatures)
}

func reduceSystemFeatureSet[T any](features *SystemFeatures, event *feature_v2.SetEvent[T]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyLoginDefaultOrg:
		features.LoginDefaultOrg.set(level, event.Value)
	case feature.KeyTriggerIntrospectionProjections:
		features.TriggerIntrospectionProjections.set(level, event.Value)
	case feature.KeyLegacyIntrospection:
		features.LegacyIntrospection.set(level, event.Value)
	case feature.KeyUserSchema:
		features.UserSchema.set(level, event.Value)
	case feature.KeyTokenExchange:
		features.TokenExchange.set(level, event.Value)
	case feature.KeyActions:
		features.Actions.set(level, event.Value)
	case feature.KeyImprovedPerformance:
		features.ImprovedPerformance.set(level, event.Value)
	case feature.KeyOIDCSingleV1SessionTermination:
		features.OIDCSingleV1SessionTermination.set(level, event.Value)
	case feature.KeyDisableUserTokenEvent:
		features.DisableUserTokenEvent.set(level, event.Value)
	case feature.KeyEnableBackChannelLogout:
		features.EnableBackChannelLogout.set(level, event.Value)
	}
	return nil
}
