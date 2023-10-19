package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
)

type SetLimits struct {
	AuditLogRetention time.Duration `json:"AuditLogRetention,omitempty"`
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
	createNew := aggregateId == ""
	if aggregateId == "" {
		aggregateId, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	cmd := c.SetLimitsCommand(ctx, limits.NewAggregate(aggregateId, instanceId, resourceOwner), wm, createNew, setLimits)
	if cmd != nil {
		events, err := c.eventstore.Push(ctx, cmd)
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
		return nil, errors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound")
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

func (c *Commands) SetLimitsCommand(ctx context.Context, a *limits.Aggregate, wm *limitsWriteModel, createNew bool, setLimits *SetLimits) eventstore.Command {
	changes := wm.NewChanges(createNew, setLimits)
	if len(changes) == 0 {
		return nil
	}
	return limits.NewSetEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			&a.Aggregate,
			limits.SetEventType,
		),
		changes...,
	)
}
