package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type InstanceFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
	UserSchema                      *bool
	TokenExchange                   *bool
	Actions                         *bool
}

func (m *InstanceFeatures) isEmpty() bool {
	return m.LoginDefaultOrg == nil &&
		m.TriggerIntrospectionProjections == nil &&
		m.LegacyIntrospection == nil &&
		m.UserSchema == nil &&
		m.TokenExchange == nil &&
		m.Actions == nil
}

func (c *Commands) SetInstanceFeatures(ctx context.Context, f *InstanceFeatures) (*domain.ObjectDetails, error) {
	if f.isEmpty() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vigh1", "Errors.NoChangesFound")
	}
	wm := NewInstanceFeaturesWriteModel(authz.GetInstance(ctx).InstanceID())
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

func prepareSetFeatures(instanceID string, f *InstanceFeatures) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			wm := NewInstanceFeaturesWriteModel(instanceID)
			return wm.setCommands(ctx, f), nil
		}, nil
	}
}

func (c *Commands) ResetInstanceFeatures(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	wm := NewInstanceFeaturesWriteModel(instanceID)
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	if wm.isEmpty() {
		return writeModelToObjectDetails(wm.WriteModel), nil
	}
	aggregate := feature_v2.NewAggregate(instanceID, instanceID)
	events, err := c.eventstore.Push(ctx, feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
