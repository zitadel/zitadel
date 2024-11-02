package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SystemFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
	TokenExchange                   *bool
	UserSchema                      *bool
	Actions                         *bool
	ImprovedPerformance             []feature.ImprovedPerformanceType
	OIDCSingleV1SessionTermination  *bool
	DisableUserTokenEvent           *bool
	EnableBackChannelLogout         *bool
}

func (m *SystemFeatures) isEmpty() bool {
	return m.LoginDefaultOrg == nil &&
		m.TriggerIntrospectionProjections == nil &&
		m.LegacyIntrospection == nil &&
		m.UserSchema == nil &&
		m.TokenExchange == nil &&
		m.Actions == nil &&
		// nil check to allow unset improvements
		m.ImprovedPerformance == nil &&
		m.OIDCSingleV1SessionTermination == nil &&
		m.DisableUserTokenEvent == nil &&
		m.EnableBackChannelLogout == nil
}

func (c *Commands) SetSystemFeatures(ctx context.Context, f *SystemFeatures) (*domain.ObjectDetails, error) {
	if f.isEmpty() {
		return nil, zerrors.ThrowInternal(nil, "COMMAND-Oop8a", "Errors.NoChangesFound")
	}
	wm := NewSystemFeaturesWriteModel()
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	cmds := wm.setCommands(ctx, f)
	if len(cmds) == 0 {
		return writeModelToObjectDetails(wm.WriteModel), nil
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) ResetSystemFeatures(ctx context.Context) (*domain.ObjectDetails, error) {
	wm := NewSystemFeaturesWriteModel()
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	if wm.isEmpty() {
		return writeModelToObjectDetails(wm.WriteModel), nil
	}
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")
	events, err := c.eventstore.Push(ctx, feature_v2.NewResetEvent(ctx, aggregate, feature_v2.SystemResetEventType))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
