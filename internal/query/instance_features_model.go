package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	feature_v1 "github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type InstanceFeaturesReadModel struct {
	*eventstore.ReadModel
	system   *SystemFeatures
	instance *InstanceFeatures
}

func NewInstanceFeaturesReadModel(ctx context.Context, system *SystemFeatures) *InstanceFeaturesReadModel {
	instanceID := authz.GetInstance(ctx).InstanceID()
	m := &InstanceFeaturesReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
		instance: new(InstanceFeatures),
		system:   system,
	}
	m.populateFromSystem()
	return m
}

func (m *InstanceFeaturesReadModel) Reduce() (err error) {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *feature_v2.ResetEvent:
			m.reduceReset()
		case *feature_v1.SetEvent[feature_v1.Boolean]:
			err = reduceInstanceFeatureSet(
				m.instance,
				feature_v1.DefaultLoginInstanceEventToV2(e),
			)
		case *feature_v2.SetEvent[bool]:
			err = reduceInstanceFeatureSet(m.instance, e)
		case *feature_v2.SetEvent[[]feature.ImprovedPerformanceType]:
			err = reduceInstanceFeatureSet(m.instance, e)
		}
		if err != nil {
			return err
		}
	}
	return m.ReadModel.Reduce()
}

func (m *InstanceFeaturesReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AddQuery().
		AggregateTypes(feature_v2.AggregateType).
		EventTypes(
			feature_v1.DefaultLoginInstanceEventType,
			feature_v2.InstanceResetEventType,
			feature_v2.InstanceLoginDefaultOrgEventType,
			feature_v2.InstanceTriggerIntrospectionProjectionsEventType,
			feature_v2.InstanceLegacyIntrospectionEventType,
			feature_v2.InstanceUserSchemaEventType,
			feature_v2.InstanceTokenExchangeEventType,
			feature_v2.InstanceActionsEventType,
			feature_v2.InstanceImprovedPerformanceEventType,
			feature_v2.InstanceWebKeyEventType,
			feature_v2.InstanceDebugOIDCParentErrorEventType,
			feature_v2.InstanceOIDCSingleV1SessionTerminationEventType,
			feature_v2.InstanceDisableUserTokenEvent,
			feature_v2.InstanceEnableBackChannelLogout,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *InstanceFeaturesReadModel) reduceReset() {
	if m.populateFromSystem() {
		return
	}
	m.instance = nil
	m.instance = new(InstanceFeatures)
}

func (m *InstanceFeaturesReadModel) populateFromSystem() bool {
	if m.system == nil {
		return false
	}
	m.instance.LoginDefaultOrg = m.system.LoginDefaultOrg
	m.instance.TriggerIntrospectionProjections = m.system.TriggerIntrospectionProjections
	m.instance.LegacyIntrospection = m.system.LegacyIntrospection
	m.instance.UserSchema = m.system.UserSchema
	m.instance.TokenExchange = m.system.TokenExchange
	m.instance.Actions = m.system.Actions
	m.instance.ImprovedPerformance = m.system.ImprovedPerformance
	m.instance.OIDCSingleV1SessionTermination = m.system.OIDCSingleV1SessionTermination
	m.instance.DisableUserTokenEvent = m.system.DisableUserTokenEvent
	m.instance.EnableBackChannelLogout = m.system.EnableBackChannelLogout
	return true
}

func reduceInstanceFeatureSet[T any](features *InstanceFeatures, event *feature_v2.SetEvent[T]) error {
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
	case feature.KeyWebKey:
		features.WebKey.set(level, event.Value)
	case feature.KeyDebugOIDCParentError:
		features.DebugOIDCParentError.set(level, event.Value)
	case feature.KeyOIDCSingleV1SessionTermination:
		features.OIDCSingleV1SessionTermination.set(level, event.Value)
	case feature.KeyDisableUserTokenEvent:
		features.DisableUserTokenEvent.set(level, event.Value)
	case feature.KeyEnableBackChannelLogout:
		features.EnableBackChannelLogout.set(level, event.Value)
	}
	return nil
}
