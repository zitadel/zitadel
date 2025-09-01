package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SystemFeatures struct {
	LoginDefaultOrg                *bool
	TokenExchange                  *bool
	UserSchema                     *bool
	ImprovedPerformance            []feature.ImprovedPerformanceType
	OIDCSingleV1SessionTermination *bool
	EnableBackChannelLogout        *bool
	LoginV2                        *feature.LoginV2
	PermissionCheckV2              *bool
}

func (m *SystemFeatures) isEmpty() bool {
	return m.LoginDefaultOrg == nil &&
		m.UserSchema == nil &&
		m.TokenExchange == nil &&
		// nil check to allow unset improvements
		m.ImprovedPerformance == nil &&
		m.OIDCSingleV1SessionTermination == nil &&
		m.EnableBackChannelLogout == nil &&
		m.LoginV2 == nil &&
		m.PermissionCheckV2 == nil
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
