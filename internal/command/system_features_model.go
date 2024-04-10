package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type SystemFeaturesWriteModel struct {
	*eventstore.WriteModel
	SystemFeatures
}

func NewSystemFeaturesWriteModel() *SystemFeaturesWriteModel {
	m := &SystemFeaturesWriteModel{
		WriteModel: &eventstore.WriteModel{
			AggregateID:   "SYSTEM",
			ResourceOwner: "SYSTEM",
		},
	}
	return m
}

func (m *SystemFeaturesWriteModel) Reduce() (err error) {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *feature_v2.ResetEvent:
			m.reduceReset()
		case *feature_v2.SetEvent[bool]:
			err = m.reduceBoolFeature(e)
		}
		if err != nil {
			return err
		}
	}
	return m.WriteModel.Reduce()
}

func (m *SystemFeaturesWriteModel) Query() *eventstore.SearchQueryBuilder {
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
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *SystemFeaturesWriteModel) reduceReset() {
	m.LoginDefaultOrg = nil
	m.TriggerIntrospectionProjections = nil
	m.LegacyIntrospection = nil
	m.TokenExchange = nil
	m.UserSchema = nil
	m.Actions = nil
}

func (m *SystemFeaturesWriteModel) reduceBoolFeature(event *feature_v2.SetEvent[bool]) error {
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
	case feature.KeyUserSchema:
		m.UserSchema = &event.Value
	case feature.KeyTokenExchange:
		m.TokenExchange = &event.Value
	case feature.KeyActions:
		m.Actions = &event.Value
	}
	return nil
}

func (wm *SystemFeaturesWriteModel) setCommands(ctx context.Context, f *SystemFeatures) []eventstore.Command {
	aggregate := feature_v2.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	cmds := make([]eventstore.Command, 0, len(feature.KeyValues())-1)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LoginDefaultOrg, f.LoginDefaultOrg, feature_v2.SystemLoginDefaultOrgEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.TriggerIntrospectionProjections, f.TriggerIntrospectionProjections, feature_v2.SystemTriggerIntrospectionProjectionsEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LegacyIntrospection, f.LegacyIntrospection, feature_v2.SystemLegacyIntrospectionEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.UserSchema, f.UserSchema, feature_v2.SystemUserSchemaEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.TokenExchange, f.TokenExchange, feature_v2.SystemTokenExchangeEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.Actions, f.Actions, feature_v2.SystemActionsEventType)
	return cmds
}

func appendFeatureUpdate[T comparable](ctx context.Context, cmds []eventstore.Command, aggregate *feature_v2.Aggregate, oldValue, newValue *T, eventType eventstore.EventType) []eventstore.Command {
	if newValue != nil && (oldValue == nil || *oldValue != *newValue) {
		cmds = append(cmds, feature_v2.NewSetEvent[T](ctx, aggregate, eventType, *newValue))
	}
	return cmds
}
