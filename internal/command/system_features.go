package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

type SystemFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
}

func (c *Commands) SetSystemFeatures(ctx context.Context, f *SystemFeatures) (*domain.ObjectDetails, error) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")
	cmds := make([]eventstore.Command, 0, len(feature.FeatureValues())-1)

	appendNonNilFeature(ctx, cmds, aggregate, f.LoginDefaultOrg, feature_v2.SystemDefaultLoginInstanceEventType)
	appendNonNilFeature(ctx, cmds, aggregate, f.TriggerIntrospectionProjections, feature_v2.SystemTriggerIntrospectionProjectionsEventType)
	appendNonNilFeature(ctx, cmds, aggregate, f.LegacyIntrospection, feature_v2.SystemLegacyIntrospectionEventType)

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) ResetSystemFeatures(ctx context.Context) (*domain.ObjectDetails, error) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")

	events, err := c.eventstore.Push(ctx, feature_v2.NewResetEvent(ctx, aggregate, feature_v2.SystemResetEventType))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
