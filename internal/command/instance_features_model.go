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
			reduceInstanceFeature(
				&m.InstanceFeatures,
				feature.KeyLoginDefaultOrg,
				feature_v1.DefaultLoginInstanceEventToV2(e).Value,
			)
		case *feature_v2.SetEvent[bool]:
			_, key, err := e.FeatureInfo()
			if err != nil {
				return err
			}
			reduceInstanceFeature(&m.InstanceFeatures, key, e.Value)
		case *feature_v2.SetEvent[*feature.LoginV2]:
			_, key, err := e.FeatureInfo()
			if err != nil {
				return err
			}
			reduceInstanceFeature(&m.InstanceFeatures, key, e.Value)
		case *feature_v2.SetEvent[[]feature.ImprovedPerformanceType]:
			_, key, err := e.FeatureInfo()
			if err != nil {
				return err
			}
			reduceInstanceFeature(&m.InstanceFeatures, key, e.Value)
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
			feature_v2.InstanceUserSchemaEventType,
			feature_v2.InstanceTokenExchangeEventType,
			feature_v2.InstanceImprovedPerformanceEventType,
			feature_v2.InstanceDebugOIDCParentErrorEventType,
			feature_v2.InstanceOIDCSingleV1SessionTerminationEventType,
			feature_v2.InstanceEnableBackChannelLogout,
			feature_v2.InstanceLoginVersion,
			feature_v2.InstancePermissionCheckV2,
			feature_v2.InstanceConsoleUseV2UserApi,
		).
		Builder().ResourceOwner(m.ResourceOwner)
}

func (m *InstanceFeaturesWriteModel) reduceReset() {
	m.InstanceFeatures = InstanceFeatures{}
}

func reduceInstanceFeature(features *InstanceFeatures, key feature.Key, value any) {
	switch key {
	case feature.KeyUnspecified:
		return
	case feature.KeyLoginDefaultOrg:
		v := value.(bool)
		features.LoginDefaultOrg = &v
	case feature.KeyTokenExchange:
		v := value.(bool)
		features.TokenExchange = &v
	case feature.KeyUserSchema:
		v := value.(bool)
		features.UserSchema = &v
	case feature.KeyImprovedPerformance:
		v := value.([]feature.ImprovedPerformanceType)
		features.ImprovedPerformance = v
	case feature.KeyDebugOIDCParentError:
		v := value.(bool)
		features.DebugOIDCParentError = &v
	case feature.KeyOIDCSingleV1SessionTermination:
		v := value.(bool)
		features.OIDCSingleV1SessionTermination = &v
	case feature.KeyEnableBackChannelLogout:
		v := value.(bool)
		features.EnableBackChannelLogout = &v
	case feature.KeyLoginV2:
		features.LoginV2 = value.(*feature.LoginV2)
	case feature.KeyPermissionCheckV2:
		v := value.(bool)
		features.PermissionCheckV2 = &v
	case feature.KeyConsoleUseV2UserApi:
		v := value.(bool)
		features.ConsoleUseV2UserApi = &v
	}
}

func (wm *InstanceFeaturesWriteModel) setCommands(ctx context.Context, f *InstanceFeatures) []eventstore.Command {
	aggregate := feature_v2.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	cmds := make([]eventstore.Command, 0, len(feature.KeyValues())-1)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LoginDefaultOrg, f.LoginDefaultOrg, feature_v2.InstanceLoginDefaultOrgEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.TokenExchange, f.TokenExchange, feature_v2.InstanceTokenExchangeEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.UserSchema, f.UserSchema, feature_v2.InstanceUserSchemaEventType)
	cmds = appendFeatureSliceUpdate(ctx, cmds, aggregate, wm.ImprovedPerformance, f.ImprovedPerformance, feature_v2.InstanceImprovedPerformanceEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.DebugOIDCParentError, f.DebugOIDCParentError, feature_v2.InstanceDebugOIDCParentErrorEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.OIDCSingleV1SessionTermination, f.OIDCSingleV1SessionTermination, feature_v2.InstanceOIDCSingleV1SessionTerminationEventType)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.EnableBackChannelLogout, f.EnableBackChannelLogout, feature_v2.InstanceEnableBackChannelLogout)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.LoginV2, f.LoginV2, feature_v2.InstanceLoginVersion)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.PermissionCheckV2, f.PermissionCheckV2, feature_v2.InstancePermissionCheckV2)
	cmds = appendFeatureUpdate(ctx, cmds, aggregate, wm.ConsoleUseV2UserApi, f.ConsoleUseV2UserApi, feature_v2.InstanceConsoleUseV2UserApi)
	return cmds
}
