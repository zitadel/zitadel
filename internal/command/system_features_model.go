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

func (m *SystemFeaturesWriteModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *feature_v2.ResetEvent:
			m.reduceReset()
		case *feature_v2.SetEvent[bool]:
			_, key, err := e.FeatureInfo()
			if err != nil {
				return err
			}
			reduceSystemFeature(&m.SystemFeatures, key, e.Value)
		case *feature_v2.SetEvent[[]feature.ImprovedPerformanceType]:
			_, key, err := e.FeatureInfo()
			if err != nil {
				return err
			}
			reduceSystemFeature(&m.SystemFeatures, key, e.Value)
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
			feature_v2.SystemImprovedPerformanceEventType,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *SystemFeaturesWriteModel) reduceReset() {
	m.SystemFeatures = SystemFeatures{}
}

func reduceSystemFeature(features *SystemFeatures, key feature.Key, value any) {
	switch key {
	case feature.KeyUnspecified:
		return
	case feature.KeyLoginDefaultOrg:
		v := value.(bool)
		features.LoginDefaultOrg = &v
	case feature.KeyTriggerIntrospectionProjections:
		v := value.(bool)
		features.TriggerIntrospectionProjections = &v
	case feature.KeyLegacyIntrospection:
		v := value.(bool)
		features.LegacyIntrospection = &v
	case feature.KeyUserSchema:
		v := value.(bool)
		features.UserSchema = &v
	case feature.KeyTokenExchange:
		v := value.(bool)
		features.TokenExchange = &v
	case feature.KeyActions:
		v := value.(bool)
		features.Actions = &v
	case feature.KeyImprovedPerformance:
		features.ImprovedPerformance = value.([]feature.ImprovedPerformanceType)
	}
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
	cmds = appendFeatureSliceUpdate(ctx, cmds, aggregate, wm.ImprovedPerformance, f.ImprovedPerformance, feature_v2.SystemImprovedPerformanceEventType)
	return cmds
}

func appendFeatureUpdate[T comparable](ctx context.Context, cmds []eventstore.Command, aggregate *feature_v2.Aggregate, oldValue, newValue *T, eventType eventstore.EventType) []eventstore.Command {
	if newValue != nil && (oldValue == nil || *oldValue != *newValue) {
		cmds = append(cmds, feature_v2.NewSetEvent[T](ctx, aggregate, eventType, *newValue))
	}
	return cmds
}

func appendFeatureSliceUpdate[T comparable](ctx context.Context, cmds []eventstore.Command, aggregate *feature_v2.Aggregate, oldValues, newValues []T, eventType eventstore.EventType) []eventstore.Command {
	if len(newValues) != len(oldValues) {
		return append(cmds, feature_v2.NewSetEvent[[]T](ctx, aggregate, eventType, newValues))
	}
	for i, oldValue := range oldValues {
		if oldValue != newValues[i] {
			return append(cmds, feature_v2.NewSetEvent[[]T](ctx, aggregate, eventType, newValues))
		}
	}
	return cmds
}
