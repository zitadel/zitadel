package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	feature_v1 "github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type InstanceFeaturesWriteModel struct {
	*eventstore.WriteModel
	InstanceFeatures
}

func NewInstanceFeaturesWriteModel(instanceID string) *InstanceFeaturesWriteModel {
	m := &InstanceFeaturesWriteModel{
		WriteModel: &eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
	}
	return m
}

func (m *InstanceFeaturesWriteModel) Reduce() (err error) {
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
	return m.WriteModel.Reduce()
}

func (m *InstanceFeaturesWriteModel) Query() *eventstore.SearchQueryBuilder {
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

func (m *InstanceFeaturesWriteModel) reduceReset() {
	m.LoginDefaultOrg = nil
	m.TriggerIntrospectionProjections = nil
	m.LegacyIntrospection = nil
	m.UserSchema = nil
	m.TokenExchange = nil
	m.Actions = nil
}

func (m *InstanceFeaturesWriteModel) reduceBoolFeature(event *feature_v2.SetEvent[bool]) error {
	_, key, err := event.FeatureInfo()
	if err != nil {
		return err
	}
	switch key {
	case feature.KeyUnspecified:
		return nil
	case feature.KeyLoginDefaultOrg:
		m.LoginDefaultOrg = &event.Value
	case feature.KeyTriggerIntrospectionProjections:
		m.TriggerIntrospectionProjections = &event.Value
	case feature.KeyLegacyIntrospection:
		m.LegacyIntrospection = &event.Value
	case feature.KeyTokenExchange:
		m.TokenExchange = &event.Value
	case feature.KeyUserSchema:
		m.UserSchema = &event.Value
	case feature.KeyActions:
		m.Actions = &event.Value
	}
	return nil
}

func (wm *InstanceFeaturesWriteModel) setCommands(ctx context.Context, f *InstanceFeatures) []eventstore.Command {
	aggregate := feature_v2.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	cmds := make([]eventstore.Command, 0, len(feature.KeyValues())-1)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LoginDefaultOrg, f.LoginDefaultOrg, feature_v2.InstanceLoginDefaultOrgEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.TriggerIntrospectionProjections, f.TriggerIntrospectionProjections, feature_v2.InstanceTriggerIntrospectionProjectionsEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LegacyIntrospection, f.LegacyIntrospection, feature_v2.InstanceLegacyIntrospectionEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.TokenExchange, f.TokenExchange, feature_v2.InstanceTokenExchangeEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.UserSchema, f.UserSchema, feature_v2.InstanceUserSchemaEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.Actions, f.Actions, feature_v2.InstanceActionsEventType)
	return cmds
}
