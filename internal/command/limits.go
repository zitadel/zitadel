package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetLimits struct {
	AuditLogRetention *time.Duration
}

// SetLimits creates new limits or updates existing limits.
func (c *Commands) SetLimits(
	ctx context.Context,
	resourceOwner string,
	setLimits *SetLimits,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getLimitsWriteModel(ctx, instanceId, resourceOwner)
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
	createCmds, err := c.SetLimitsCommand(limits.NewAggregate(aggregateId, instanceId, resourceOwner), wm, setLimits)()
	if err != nil {
		return nil, err
	}
	cmds, err := createCmds(ctx, nil)
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

func (c *Commands) ResetLimits(ctx context.Context, resourceOwner string) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getLimitsWriteModel(ctx, instanceId, resourceOwner)
	if err != nil {
		return nil, err
	}
	if wm.AggregateID == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound")
	}
	aggregate := limits.NewAggregate(wm.AggregateID, instanceId, resourceOwner)
	events := []eventstore.Command{limits.NewResetEvent(ctx, &aggregate.Aggregate)}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getLimitsWriteModel(ctx context.Context, instanceId, resourceOwner string) (*limitsWriteModel, error) {
	wm := newLimitsWriteModel(instanceId, resourceOwner)
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
}

func (c *Commands) SetLimitsCommand(a *limits.Aggregate, wm *limitsWriteModel, setLimits *SetLimits) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if setLimits == nil || setLimits.AuditLogRetention == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4M9vs", "Errors.Limits.NoneSpecified")
		}
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			changes := wm.NewChanges(setLimits)
			if len(changes) == 0 {
				return nil, nil
			}
			return []eventstore.Command{limits.NewSetEvent(
				eventstore.NewBaseEventForPush(
					ctx,
					&a.Aggregate,
					limits.SetEventType,
				),
				changes...,
			)}, nil
		}, nil
	}
}
