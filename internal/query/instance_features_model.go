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
			err = m.reduceBoolFeature(
				feature_v1.DefaultLoginInstanceEventToV2(e),
			)
		case *feature_v2.SetEvent[bool]:
			err = m.reduceBoolFeature(e)
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
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *InstanceFeaturesReadModel) reduceReset() {
	if m.populateFromSystem() {
		return
	}
	m.instance.LoginDefaultOrg = FeatureSource[bool]{}
	m.instance.TriggerIntrospectionProjections = FeatureSource[bool]{}
	m.instance.LegacyIntrospection = FeatureSource[bool]{}
	m.instance.UserSchema = FeatureSource[bool]{}
	m.instance.TokenExchange = FeatureSource[bool]{}
	m.instance.Actions = FeatureSource[bool]{}
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
	return true
}

func (m *InstanceFeaturesReadModel) reduceBoolFeature(event *feature_v2.SetEvent[bool]) error {
	level, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	var dst *FeatureSource[bool]

	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyLoginDefaultOrg:
		dst = &m.instance.LoginDefaultOrg
	case feature.KeyTriggerIntrospectionProjections:
		dst = &m.instance.TriggerIntrospectionProjections
	case feature.KeyLegacyIntrospection:
		dst = &m.instance.LegacyIntrospection
	case feature.KeyUserSchema:
		dst = &m.instance.UserSchema
	case feature.KeyTokenExchange:
		dst = &m.instance.TokenExchange
	case feature.KeyActions:
		dst = &m.instance.Actions
	}
	*dst = FeatureSource[bool]{
		Level: level,
		Value: event.Value,
	}
	return nil
}
