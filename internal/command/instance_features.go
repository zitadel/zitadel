package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type InstanceFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
}

func appendNonNilFeature[T any](ctx context.Context, cmds []eventstore.Command, aggregate *feature_v2.Aggregate, value *T, eventType eventstore.EventType) []eventstore.Command {
	if value != nil {
		cmds = append(cmds, feature_v2.NewSetEvent[T](ctx, aggregate, eventType, *value))
	}
	return cmds
}

func (c *Commands) SetInstanceFeatures(ctx context.Context, f *InstanceFeatures) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	aggregate := feature_v2.NewAggregate(instanceID, instanceID)
	cmds := make([]eventstore.Command, 0, len(feature.FeatureValues())-1)

	appendNonNilFeature(ctx, cmds, aggregate, f.LoginDefaultOrg, feature_v2.InstanceDefaultLoginInstanceEventType)
	appendNonNilFeature(ctx, cmds, aggregate, f.TriggerIntrospectionProjections, feature_v2.InstanceTriggerIntrospectionProjectionsEventType)
	appendNonNilFeature(ctx, cmds, aggregate, f.LegacyIntrospection, feature_v2.InstanceLegacyIntrospectionEventType)

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) ResetInstanceFeatures(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	aggregate := feature_v2.NewAggregate(instanceID, instanceID)

	events, err := c.eventstore.Push(ctx, feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
