package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

type SetRestrictions struct {
	DisallowPublicOrgRegistration *bool
}

// SetRestrictions creates new restrictions or updates existing restrictions.
func (c *Commands) SetInstanceRestrictions(
	ctx context.Context,
	setRestrictions *SetRestrictions,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getRestrictionsWriteModel(ctx, instanceId, instanceId)
	if err != nil {
		return nil, err
	}
	aggregateId := wm.AggregateID
	if aggregateId == "" {
		aggregateId, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}
	setCmd, err := c.SetRestrictionsCommand(restrictions.NewAggregate(aggregateId, instanceId, instanceId), wm, setRestrictions)()
	if err != nil {
		return nil, err
	}
	cmds, err := setCmd(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(cmds) > 0 {
		events, err := c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(wm, events...)
		if err != nil {
			return nil, err
		}
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getRestrictionsWriteModel(ctx context.Context, instanceId, resourceOwner string) (*restrictionsWriteModel, error) {
	wm := newRestrictionsWriteModel(instanceId, resourceOwner)
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
}

func (c *Commands) SetRestrictionsCommand(a *restrictions.Aggregate, wm *restrictionsWriteModel, setRestrictions *SetRestrictions) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if setRestrictions == nil || setRestrictions.DisallowPublicOrgRegistration == nil {
			return nil, errors.ThrowInvalidArgument(nil, "COMMAND-oASwj", "Errors.Restrictions.NoneSpecified")
		}
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			changes := wm.NewChanges(setRestrictions)
			if len(changes) == 0 {
				return nil, nil
			}
			return []eventstore.Command{restrictions.NewSetEvent(
				eventstore.NewBaseEventForPush(
					ctx,
					&a.Aggregate,
					restrictions.SetEventType,
				),
				changes...,
			)}, nil
		}, nil
	}
}
