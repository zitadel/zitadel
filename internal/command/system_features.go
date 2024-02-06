package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SystemFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
}

func (c *Commands) SetSystemFeatures(ctx context.Context, f *SystemFeatures) (*domain.ObjectDetails, error) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")
	cmds := make([]eventstore.Command, 0, len(feature.KeyValues())-1)

	cmds = appendNonNilFeature(ctx, cmds, aggregate, f.LoginDefaultOrg, feature_v2.SystemDefaultLoginInstanceEventType)
	cmds = appendNonNilFeature(ctx, cmds, aggregate, f.TriggerIntrospectionProjections, feature_v2.SystemTriggerIntrospectionProjectionsEventType)
	cmds = appendNonNilFeature(ctx, cmds, aggregate, f.LegacyIntrospection, feature_v2.SystemLegacyIntrospectionEventType)
	if len(cmds) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Dul8I", "Errors.NoChangesFound")
	}
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
